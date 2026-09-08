# AGENTS.md

This file provides guidance to automated agents (AI or otherwise) when working with code in this repository.

## Project Overview

The eventrouter is a Kubernetes event watcher that streams `v1.Event` resources to configurable sinks. It runs as a deployment inside a cluster, watches for event creates/updates/deletes via a shared informer, and serializes them as JSON to the chosen sink. In OpenShift Logging, it feeds Kubernetes events into the log collection pipeline so they can be stored and queried alongside container logs.

## Architecture

### Core Components

1. **Main Entry Point** (`main.go`)
   - Loads config from JSON file at `/etc/eventrouter/config` (or current directory)
   - Builds a Kubernetes clientset (in-cluster or via kubeconfig)
   - Sets up a shared informer for `v1.Events` with configurable resync interval
   - Optionally exposes Prometheus metrics on `:8080/metrics` and pprof on `/debug/pprof/`
   - Handles graceful shutdown via OS signals

2. **EventRouter** (`eventrouter.go`)
   - Registers add/update/delete event handlers on the events informer
   - On add: wraps event as `EventData{Verb: "ADDED"}` and sends to sink
   - On update: skips if ResourceVersion unchanged, otherwise wraps as `EventData{Verb: "UPDATED"}` with old event
   - On delete: logs only (TTL expiration), does not forward to sink
   - Runs until stop channel is closed

3. **Sinks** (`sinks/`)
   - `EventSinkInterface`: single method `UpdateEvents(eNew, eOld *v1.Event)`
   - `ManufactureSink()`: factory that reads `sink` from viper config
   - Available sinks:
     - `glog` (default) — JSON via glog, useful with existing EFK stacks
     - `stdout` — raw JSON to stdout for direct Fluentd/ES indexing
     - `http` — buffered HTTP POST with configurable buffer size and overflow policy
     - `kafka` — Kafka producer with configurable brokers, topic, async mode

4. **EventData** (`sinks/eventdata.go`)
   - Wraps new/old event with a verb field (`ADDED`/`UPDATED`)
   - Supports JSON serialization and RFC 5424 syslog format

### Configuration

Config is a JSON file with these keys:

| Key | Default | Description |
|-----|---------|-------------|
| `sink` | `glog` | Sink type: `glog`, `stdout`, `http`, `kafka` |
| `kubeconfig` | `""` | Path to kubeconfig; empty uses in-cluster config |
| `resync-interval` | `30m` | Informer resync interval |
| `enable-prometheus` | `true` | Expose Prometheus metrics |
| `enable-http-pprof` | `false` | Expose pprof endpoints |
| `WATCH_NAMESPACE` | `""` | Restrict to a single namespace (env var) |
| `httpSinkUrl` | — | Required for `http` sink |
| `httpSinkBufferSize` | `1500` | Event buffer size for `http` sink |
| `httpSinkDiscardMessages` | `true` | Drop events on buffer overflow |
| `kafkaBrokers` | `["kafka:9092"]` | Kafka broker addresses |
| `kafkaTopic` | `eventrouter` | Kafka topic name |
| `kafkaAsync` | `true` | Use async Kafka producer |
| `kafkaRetryMax` | `5` | Max Kafka send retries |

## Development Workflow

### Building
```bash
make build    # Builds the eventrouter binary
make fmt      # Runs gofmt
make image    # Builds container image with podman
```

### Testing
```bash
make test                        # Runs tests in sinks/...
go test -v ./sinks/...           # Verbose test output
```

### Deployment
```bash
kubectl create -f yaml/eventrouter.yaml              # Cluster-wide
kubectl create -f yaml/eventrouter-namespaced.yaml    # Single namespace
```

## Code Conventions

- Go standard style enforced by `gofmt`
- Module path: `github.com/openshift/eventrouter`
- Logging via `github.com/golang/glog`
- Config via `github.com/spf13/viper`
- Kubernetes client-go informer pattern for event watching
- Fork-based git workflow: push to your fork, PR against `origin/master`

## References

- **Architecture**: [ARCHITECTURE.md](ARCHITECTURE.md)
- **Contributing**: [CONTRIBUTING.md](CONTRIBUTING.md)
- **README**: [README.md](README.md)

## Key Dependencies

- `k8s.io/client-go` — Kubernetes client, informers, listers
- `github.com/spf13/viper` — Configuration
- `github.com/golang/glog` — Structured logging
- `github.com/prometheus/client_golang` — Prometheus metrics
- `github.com/IBM/sarama` — Kafka client
- `github.com/sethgrid/pester` — Resilient HTTP client (used by HTTP sink)
- `github.com/crewjam/rfc5424` — Syslog message formatting
