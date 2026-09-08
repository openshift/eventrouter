# Architecture: eventrouter

This document describes the internal design, key decisions, and implementation details.

## Overview

The eventrouter watches the Kubernetes Event API via a shared informer and forwards each event as a JSON-serialized message to a configurable sink. It runs as a single-pod deployment inside a cluster — no leader election, no HA, no persistent state. In OpenShift Logging, it feeds Kubernetes events into the log collection pipeline (typically via stdout → Vector → LokiStack) so they can be stored and queried alongside container logs.

**Main output**: JSON objects wrapping `v1.Event` with a verb (`ADDED` or `UPDATED`) and optionally the previous event.

## System Architecture

```text
  KUBERNETES API                       SINKS

┌─────────────────────┐
│  kube-apiserver      │
│  (Event resources)   │
└──────────┬──────────┘
           │ watch/list
           │
┌──────────▼──────────┐
│  SharedInformer      │
│  (client-go)         │
│  • Resync every 30m  │
│  • Namespace filter  │
│    (optional)        │
└──────────┬──────────┘
           │ Add/Update/Delete callbacks
           │
┌──────────▼──────────┐          ┌──────────────────────┐
│  EventRouter         │          │  GlogSink            │
│  (eventrouter.go)    │──sink───▶│  StdoutSink          │
│  • Wraps events as   │          │  HTTPSink            │
│    EventData JSON    │          │  KafkaSink           │
│  • Skips no-op       │          └──────────────────────┘
│    updates           │
│  • Logs deletes only │
└─────────────────────┘

  OPTIONAL ENDPOINTS (main.go)

┌─────────────────────┐
│  HTTP server :8080   │
│  • /metrics          │
│    (Prometheus)      │
│  • /debug/pprof/*    │
│    (optional)        │
└─────────────────────┘
```

## Component Details

### 1. Main Entry Point (`main.go`)

**Responsibilities**:
- Load JSON config from `/etc/eventrouter/config` or current directory via viper
- Build Kubernetes clientset (in-cluster or via `kubeconfig` / `KUBECONFIG` env var)
- Create a `SharedInformerFactory`, optionally scoped to `WATCH_NAMESPACE`
- Wire up `EventRouter` with the events informer
- Optionally start HTTP server for Prometheus metrics and/or pprof
- Block on OS signal (SIGINT, SIGTERM, etc.) for graceful shutdown

**Key Design Decisions**:
- **Viper for config**: JSON file + env var overrides; no CLI flags for config values (only `-listen-address` and glog flags)
- **No TLS on metrics endpoint**: Plain HTTP via `ListenAndServe` — TLS is handled by the service mesh or sidecar in OpenShift
- **Exit code 1 on shutdown**: `os.Exit(1)` after WaitGroup completes, even on clean shutdown — Kubernetes restarts the pod regardless

**Concurrency**:
- Main goroutine blocks on `wg.Wait()` after starting the informer
- EventRouter goroutine runs `eventRouter.Run(stop)`, which blocks until the stop channel closes
- HTTP server (if enabled) runs in its own goroutine
- Signal handler goroutine closes the stop channel on signal

### 2. EventRouter (`eventrouter.go`)

**Responsibilities**:
- Register add/update/delete handlers on the events informer
- Forward new/updated events to the configured sink
- Skip updates where `ResourceVersion` is unchanged (no-op resync)
- Log deletions (TTL expiration) without forwarding

**Key Design Decisions**:
- **Single sink**: The router supports exactly one sink at a time, selected by config. Multiple sinks were considered (see TODO in code) but not implemented.
- **No filtering**: All events in the watched namespace(s) are forwarded. Filtering by reason, type, or involved object is left to downstream consumers.
- **Delete = log only**: Event deletions are TTL-driven garbage collection, not meaningful user events, so they are logged at V(5) but not sent to the sink.

**Event Flow**:
| Informer Callback | Behavior |
|-------------------|----------|
| `addEvent` | Wrap as `EventData{Verb: "ADDED", Event: new}` → sink |
| `updateEvent` | Skip if `ResourceVersion` unchanged; else wrap as `EventData{Verb: "UPDATED", Event: new, OldEvent: old}` → sink |
| `deleteEvent` | Log at V(5) only, do not forward |

### 3. Sinks (`sinks/`)

All sinks implement `EventSinkInterface`:

```go
type EventSinkInterface interface {
    UpdateEvents(eNew *v1.Event, eOld *v1.Event)
}
```

`ManufactureSink()` reads the `sink` config key and constructs the appropriate implementation.

#### GlogSink (`glogsink.go`)
- JSON-serializes `EventData` and logs via `glog.Info`
- Default sink; useful with existing EFK stacks that index glog output
- No buffering, no goroutines — synchronous in the informer callback

#### StdoutSink (`stdoutsink.go`)
- JSON-serializes `EventData` and writes to stdout via `fmt.Println`
- Preferred when log collectors (Fluentd, Vector) read container stdout directly
- No glog formatting overhead — raw JSON lines

#### HTTPSink (`httpsink.go`)
- Sends events as RFC 5424 syslog messages over HTTP POST
- Compatible with Heroku Logplex HTTP drain protocol
- Asynchronous: events are written to a buffered channel, a background goroutine drains and batches them into HTTP requests
- Uses `pester` client with exponential jitter backoff (max 10 retries)
- Configurable buffer size (default 1500) and overflow policy (discard or block)
- Reuses a `bytes.Buffer` for request bodies to minimize allocations
- Coalesces multiple buffered events into a single HTTP request

#### KafkaSink (`kafkasink.go`)
- Sends JSON-serialized events to a Kafka topic
- Uses `sarama` client library
- Supports sync and async producers (config: `kafkaAsync`)
- Message key is the involved object name (`eNew.InvolvedObject.Name`)
- Sync mode: `WaitForAll` acks, returns partition/offset
- Async mode: non-blocking send to input channel, logs errors from error channel

### 4. EventData (`sinks/eventdata.go`)

- Wraps `v1.Event` with verb and optional old event
- Supports JSON marshaling (used by all sinks)
- Supports RFC 5424 syslog format via `WriteRFC5424()` (used by HTTP sink)
  - Uses `crewjam/rfc5424` library
  - Hostname from `Event.Source.Host`, AppName from `Event.Source.Component`
  - Message body is the JSON-serialized event data

## Concurrency Model

The eventrouter has minimal concurrency:

1. **Informer goroutines** (managed by client-go) deliver events to the registered callbacks
2. **Callbacks run synchronously** in the informer's event processing goroutine — this means sink operations (glog, stdout, channel write) must not block for long
3. **HTTPSink** is the exception: `UpdateEvents` writes to a channel (non-blocking with overflow), and a separate goroutine (`Run`) drains the channel and makes HTTP calls
4. **KafkaSink async mode**: `UpdateEvents` writes to sarama's input channel, sarama manages its own goroutines

There is no explicit mutex in the eventrouter itself — thread safety comes from the informer's sequential event delivery and the sinks' internal synchronization.

## Configuration

Config is loaded from a JSON file (`/etc/eventrouter/config` or `./config.json`):

```json
{
  "sink": "glog",
  "kubeconfig": "",
  "resync-interval": "30m",
  "enable-prometheus": true,
  "enable-http-pprof": false
}
```

Environment variable overrides:
- `KUBECONFIG` — path to kubeconfig file
- `WATCH_NAMESPACE` — restrict event watching to a single namespace
- `EVENTROUTER_CONFIG` — override config file path

### Sink-Specific Config

**HTTP sink**:
```json
{
  "sink": "http",
  "httpSinkUrl": "http://logplex:8080/events",
  "httpSinkBufferSize": 1500,
  "httpSinkDiscardMessages": true
}
```

**Kafka sink**:
```json
{
  "sink": "kafka",
  "kafkaBrokers": ["kafka:9092"],
  "kafkaTopic": "eventrouter",
  "kafkaAsync": true,
  "kafkaRetryMax": 5
}
```

## Container Image

Multi-stage build (`Dockerfile`):
1. **Builder**: `golang:1.26.7` — compiles the binary
2. **Runtime**: `ubi9/ubi-minimal` — runs as UID 1000 (non-root)
3. **Entrypoint**: `/bin/eventrouter -v 3 -logtostderr`

The `-v 3` flag sets glog verbosity to 3 by default, and `-logtostderr` sends glog output to stderr (which container runtimes capture as the pod's log stream).

## Deployment

Two deployment manifests in `yaml/`:
- `eventrouter.yaml` — cluster-wide, watches all namespaces
- `eventrouter-namespaced.yaml` — single namespace via `WATCH_NAMESPACE`

Both create:
- A `ServiceAccount`
- A `ClusterRole` / `ClusterRoleBinding` (or `Role` / `RoleBinding`) with `get`, `list`, `watch` on events
- A `ConfigMap` with the JSON config
- A single-replica `Deployment`

## Testing Strategy

### Unit Tests
- `sinks/httpsink_test.go` — HTTP sink tests

### Running Tests
```bash
make test                        # Runs sinks/... tests
go test -v ./sinks/...           # Verbose output
```

### Test Coverage
Test coverage is limited — only the HTTP sink has unit tests. The glog, stdout, and Kafka sinks have no tests. The EventRouter itself has no unit tests; correctness relies on the client-go informer contract.

## Known Tradeoffs

### 1. Single Sink vs. Multiple
**Decision**: One sink at a time.
- **Pro**: Simple config, simple code, no fan-out complexity
- **Con**: Can't simultaneously write to stdout and Kafka
- **Rationale**: In OpenShift Logging, the collection pipeline (Vector) handles fan-out. The eventrouter just needs to get events into the pipeline.

### 2. Synchronous Callbacks (glog/stdout) vs. Async (HTTP/Kafka)
**Decision**: glog and stdout sinks run synchronously in the informer callback.
- **Pro**: Zero buffering overhead, guaranteed delivery order
- **Con**: A slow sink could back-pressure the informer
- **Rationale**: glog and stdout are local I/O — effectively non-blocking. HTTP and Kafka are remote, so they use async patterns.

### 3. No Event Filtering
**Decision**: Forward all events.
- **Pro**: No config complexity, no risk of dropping important events
- **Con**: High-volume clusters generate many events; downstream must filter
- **Rationale**: Filtering is better done in the log pipeline (Vector/CLF) where it can be configured declaratively.

### 4. No Leader Election / HA
**Decision**: Single-replica deployment, no leader election.
- **Pro**: Simple operations, no coordination overhead
- **Con**: Brief gap in event coverage during pod restarts
- **Rationale**: Events have a TTL in etcd (default 1 hour). On restart, the informer does a full list, so events are not permanently lost — only the real-time stream is interrupted.

### 5. Exit Code 1 on Clean Shutdown
**Decision**: `os.Exit(1)` after the WaitGroup completes, even on a clean signal-driven shutdown.
- **Pro**: Kubernetes always restarts the pod (desired behavior for a long-running watcher)
- **Con**: Misleading exit code — 1 typically indicates an error
- **Rationale**: Legacy behavior; changing it would have no practical effect since the pod is managed by a Deployment with restart policy Always.
