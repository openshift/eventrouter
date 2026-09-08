# Architecture: eventrouter

This document describes the internal design, key decisions, and implementation details.

## Overview

The eventrouter watches the Kubernetes Event API via a shared informer and forwards each created or updated event as a JSON-serialized message to stdout (deletes are logged, not forwarded). It runs as a single-pod deployment inside a cluster — no leader election, no HA, no persistent state. In OpenShift Logging, it feeds Kubernetes events into the log collection pipeline (stdout → Vector → LokiStack) so they can be stored and queried alongside container logs.

**Main output**: JSON objects wrapping `v1.Event` with a verb (`ADDED` or `UPDATED`) and optionally the previous event.

**Sink model**: The stdout sink is the only sink and is hardcoded (`sinks.ManufactureSink()` returns it unconditionally). The glog, HTTP, and Kafka sinks were never supported and have been removed to shrink the dependency and CVE surface ([LOG-10393](https://redhat.atlassian.net/browse/LOG-10393)); a `make check-deps` guard keeps their dependency families out of the module. Forwarding to other systems is delegated to the log collector that reads container stdout.

## System Architecture

```text
  KUBERNETES API                       OUTPUT

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
│  EventRouter         │          │  StdoutSink          │
│  (eventrouter.go)    │──sink───▶│  (only sink,         │
│  • Wraps events as   │          │   hardcoded)         │
│    EventData JSON    │          │  JSON lines on       │
│  • Skips no-op       │          │  stdout              │
│    updates           │          └──────────┬───────────┘
│  • Logs deletes only │                     │
└─────────────────────┘                     ▼
                                    log collector (Vector)
┌─────────────────────┐
│  HTTP server :8080   │
│  • /metrics          │
│    (Prometheus,      │
│     default on)      │
└─────────────────────┘
```

## Component Details

### 1. Main Entry Point (`main.go`)

**Responsibilities**:
- Load JSON config via viper: search `/etc/eventrouter/` then the current directory for `config` with a supported extension (`.json` first; extensionless `config` is a fallback only when no `config.<ext>` exists; the manifests mount `/etc/eventrouter/config.json`). A config file is required — startup panics if none is found
- Build Kubernetes clientset (in-cluster or via `kubeconfig` / `KUBECONFIG` env var)
- Create a `SharedInformerFactory`, optionally scoped to `WATCH_NAMESPACE`
- Wire up `EventRouter` with the events informer
- Start the Prometheus metrics HTTP server (`enable-prometheus`, on by default)
- Block on OS signal (SIGINT, SIGTERM, etc.) for shutdown

**Key Design Decisions**:
- **Viper for config**: JSON file; only `KUBECONFIG` and `WATCH_NAMESPACE` override config values through the environment (no `AutomaticEnv` — other env vars do not map to config keys); no CLI flags for config values (only `-listen-address` and glog flags)
- **No TLS on metrics endpoint**: Plain HTTP via `ListenAndServe` — TLS is a deployment requirement, not a guarantee. Deployments must place this endpoint behind an authenticated TLS boundary (e.g., service mesh or sidecar); the manifests in this repo do not provide one
- **Exit code 1 on shutdown**: `os.Exit(1)` after WaitGroup completes, even on clean shutdown — Kubernetes restarts the pod regardless

**Concurrency**:
- Main goroutine blocks on `wg.Wait()` after starting the informer
- EventRouter goroutine runs `eventRouter.Run(stop)`, which blocks until the stop channel closes
- HTTP server (if enabled) runs in its own goroutine
- Signal handler goroutine closes the stop channel on signal

### 2. EventRouter (`eventrouter.go`)

**Responsibilities**:
- Register add/update/delete handlers on the events informer
- Forward new/updated events to the sink
- Skip updates where `ResourceVersion` is unchanged (no-op resync)
- Log deletions (TTL expiration) without forwarding

**Key Design Decisions**:
- **Hardcoded sink**: The router uses the stdout sink returned by `ManufactureSink()`; there is no runtime sink selection. Multiple sinks remain a TODO in the code.
- **No filtering**: All events in the watched namespace(s) are forwarded. Filtering by reason, type, or involved object is left to downstream consumers.
- **Delete = log only**: Event deletions are TTL-driven garbage collection, not meaningful user events, so they are logged at V(5) but not sent to the sink.

**Event Flow**:
| Informer Callback | Behavior |
|-------------------|----------|
| `addEvent` | Wrap as `EventData{Verb: "ADDED", Event: new}` → sink |
| `updateEvent` | Skip if `ResourceVersion` unchanged; else wrap as `EventData{Verb: "UPDATED", Event: new, OldEvent: old}` → sink |
| `deleteEvent` | Log at V(5) only, do not forward |

### 3. Sinks (`sinks/`)

The sink interface:

```go
type EventSinkInterface interface {
    UpdateEvents(eNew *v1.Event, eOld *v1.Event)
}
```

`ManufactureSink()` returns the stdout sink unconditionally; a stale `sink` key in an existing configuration is ignored (covered by `sinks/interfaces_test.go`).

#### StdoutSink (`stdoutsink.go`)
- JSON-serializes `EventData` and writes to stdout via `fmt.Println`
- The only compiled-in sink; log collectors (Fluentd, Vector) read container stdout directly
- No glog formatting overhead — raw JSON lines
- No buffering or background goroutine — synchronous in the informer callback (a channel/buffer remains a TODO in the code)
- Serialization errors are reported to stderr, not stdout, to keep the stdout stream clean JSON

### 4. EventData (`sinks/eventdata.go`)

- Wraps `v1.Event` with verb and optional old event
- Supports JSON marshaling (used by the stdout sink)
- RFC 5424 syslog formatting was removed together with the HTTP sink

## Concurrency Model

The eventrouter has minimal concurrency:

1. **Informer goroutines** (managed by client-go) deliver events to the registered callbacks
2. **Callbacks run synchronously** in the informer's event processing goroutine — the stdout write happens inline and must not block for long
3. **HTTP server** (optional) runs in its own goroutine and is independent of event processing

There is no explicit mutex in the eventrouter itself — thread safety comes from the informer's sequential event delivery.

## Configuration

Config is loaded by viper, which searches `/etc/eventrouter/` then the current directory for `config` with a supported extension (`.json` first; extensionless `config` is a fallback only when no `config.<ext>` exists — the manifests mount the ConfigMap as `/etc/eventrouter/config.json`). A config file is required; startup panics if none is found. The manifests ship an empty `{}` so all defaults apply:

```json
{
  "kubeconfig": "",
  "resync-interval": "30m",
  "enable-prometheus": true
}
```

Environment variable overrides:
- `KUBECONFIG` — path to kubeconfig file
- `WATCH_NAMESPACE` — restrict event watching to a single namespace

## Container Image

Multi-stage build (`Dockerfile`):
1. **Builder**: `golang:1.26.7` — compiles the binary via `make build` (which also runs the dependency guard)
2. **Runtime**: `ubi9/ubi-micro` — runs as UID 1000 (non-root)
3. **Command**: `/bin/eventrouter -v 3 -logtostderr`

The `-v 3` flag sets glog verbosity to 3 by default, and `-logtostderr` sends glog output to stderr while event JSON goes to stdout. A default pod log request returns both streams interleaved, so consumers that expect only JSON must filter by stream or tolerate glog lines.

## Deployment

Two deployment manifests in `yaml/`:
- `eventrouter.yaml` — cluster-wide, watches all namespaces
- `eventrouter-namespaced.yaml` — single namespace via `WATCH_NAMESPACE`

Both create:
- A `ServiceAccount`
- A `ClusterRole` / `ClusterRoleBinding` with `get`, `list`, `watch` on events — both manifests grant cluster-scoped read on events; the namespaced variant narrows only what the informer watches (`WATCH_NAMESPACE`), not the RBAC grant
- A `ConfigMap` with the JSON config
- A single-replica `Deployment`

## Testing Strategy

### Unit Tests
- `sinks/stdoutsink_test.go` — stdout sink serialization for ADDED/UPDATED events, stdout capture with restore-on-panic
- `sinks/interfaces_test.go` — `ManufactureSink()` returns the stdout sink and ignores a stale `sink` config key

### Running Tests
```bash
make test                        # Runs sinks/... tests (bracketed by check-deps)
go test -v ./sinks/...           # Verbose output
```

### Test Coverage
The main package (`main.go`, `eventrouter.go`) has no unit tests; correctness relies on the client-go informer contract.

## Known Tradeoffs

### 1. Stdout-Only Sink
**Decision**: Only the stdout sink is compiled in; glog, HTTP, and Kafka sinks were removed.
- **Pro**: Minimal dependency footprint and CVE surface, enforced by `make check-deps`
- **Con**: No direct delivery to Kafka/HTTP/syslog from the eventrouter itself
- **Rationale**: The removed sinks were never supported. In OpenShift Logging, the collection pipeline (Vector) handles fan-out to arbitrary destinations; the eventrouter just needs to get events into the pipeline ([LOG-10393](https://redhat.atlassian.net/browse/LOG-10393)).

### 2. Synchronous Callbacks
**Decision**: The stdout sink runs synchronously in the informer callback.
- **Pro**: Zero buffering overhead, guaranteed delivery order
- **Con**: A stalled stdout write blocks the handler goroutine — client-go queues subsequent notifications in an unbounded buffer, so memory grows without bound until the process is OOM-killed (the sink has no buffering or timeout; a channel/buffer is a TODO in the code)
- **Rationale**: stdout is local I/O to the container runtime — effectively non-blocking.

### 3. No Event Filtering
**Decision**: Forward all events.
- **Pro**: No config complexity, and filtering itself never discards events (a failed stdout write can still silently lose one — `fmt.Println` errors are ignored; only serialization errors are logged)
- **Con**: High-volume clusters generate many events; downstream must filter
- **Rationale**: Filtering is better done in the log pipeline (Vector/CLF) where it can be configured declaratively.

### 4. No Leader Election / HA
**Decision**: Single-replica deployment, no leader election.
- **Pro**: Simple operations, no coordination overhead
- **Con**: Brief gap in event coverage during pod restarts
- **Rationale**: Events have a TTL in etcd (default 1 hour). On restart, the informer does a full list and forwards the current Event objects as new events — intermediate updates made during downtime are not replayed, so their earlier values are lost even for restarts shorter than the TTL. If downtime exceeds the event TTL, expired events are permanently missed.

### 5. Exit Code 1 on Clean Shutdown
**Decision**: `os.Exit(1)` after the WaitGroup completes, even on a clean signal-driven shutdown.
- **Pro**: Kubernetes always restarts the pod (desired behavior for a long-running watcher)
- **Con**: Misleading exit code — 1 typically indicates an error
- **Rationale**: Legacy behavior; changing it would have no practical effect since the pod is managed by a Deployment with restart policy Always.
