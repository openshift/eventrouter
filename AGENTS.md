# AGENTS.md

This file provides guidance to automated agents (AI or otherwise) when working with code in this repository.

## Project Overview

The eventrouter is a Kubernetes event watcher that streams `v1.Event` resources to stdout as JSON. It runs as a deployment inside a cluster, watches the events API via a shared informer, and serializes event creates and updates as JSON lines on stdout (deletes are logged at high verbosity, not forwarded). In OpenShift Logging, it feeds Kubernetes events into the log collection pipeline so they can be stored and queried alongside container logs. Only the `stdout` sink is compiled in; forwarding events to other systems (Kafka, HTTP, syslog, ...) is the job of a log collector reading container stdout.

## Architecture

### Core Components

1. **Main Entry Point** (`main.go`)
   - Loads config via viper, which searches `/etc/eventrouter/` then the current directory for `config` with a supported extension (`.json` first; extensionless `config` is a fallback only when no `config.<ext>` exists); a config file is required — startup panics without one
   - Builds a Kubernetes clientset (in-cluster or via kubeconfig)
   - Sets up a shared informer for `v1.Events` with configurable resync interval
   - Exposes Prometheus metrics on `:8080/metrics` (enabled by default; disable with `enable-prometheus: false`)
   - Handles shutdown via OS signals

2. **EventRouter** (`eventrouter.go`)
   - Registers add/update/delete event handlers on the events informer
   - On add: wraps event as `EventData{Verb: "ADDED"}` and sends to sink
   - On update: skips if ResourceVersion unchanged, otherwise wraps as `EventData{Verb: "UPDATED"}` with old event
   - On delete: logs only (TTL expiration), does not forward to sink
   - Runs until stop channel is closed

3. **Sinks** (`sinks/`)
   - `EventSinkInterface`: single method `UpdateEvents(eNew, eOld *v1.Event)`
   - `ManufactureSink()`: returns the stdout sink — there is no runtime sink selection; a stale `sink` key in an existing config is ignored
   - `stdout` (the only sink) — raw JSON lines to stdout for direct collector indexing

4. **EventData** (`sinks/eventdata.go`)
   - Wraps new/old event with a verb field (`ADDED`/`UPDATED`)
   - Supports JSON serialization only (RFC 5424 syslog support was removed along with the HTTP sink)

### Configuration

Config is a JSON file with these keys:

| Key | Default | Description |
|-----|---------|-------------|
| `kubeconfig` | `""` | Path to kubeconfig; empty uses in-cluster config |
| `resync-interval` | `30m` | Informer resync interval |
| `enable-prometheus` | `true` | Expose Prometheus metrics |
| `WATCH_NAMESPACE` | `""` | Restrict to a single namespace (env var) |

Environment variable overrides: `KUBECONFIG` (kubeconfig path), `WATCH_NAMESPACE` (namespace filter). These are the only config keys overridable via environment.

Historical `sink`, `httpSink*`, and `kafka*` keys are no longer read.

## Development Workflow

### Building
```bash
make build    # fmt + check-deps, go build, check-deps again
make fmt      # Placeholder — only prints "gofmt"; run `gofmt -l -w .` to format
make image    # Builds container image with podman
```

### Testing
```bash
make test                        # check-deps + tests in sinks/...
go test -v ./sinks/...           # Verbose test output
```

### Dependency guard
`make check-deps` fails the build if any removed sink dependency family (sarama, pester, rfc5424, and related) reappears in `go.mod` or `go.sum`. Do not reintroduce these dependencies ([LOG-10393](https://redhat.atlassian.net/browse/LOG-10393)).

### Deployment
```bash
kubectl create -f yaml/eventrouter.yaml              # Cluster-wide
kubectl create -f yaml/eventrouter-namespaced.yaml    # Single namespace
```

## Code Conventions

- Go standard style — follow `gofmt` formatting (not enforced by CI)
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
- `github.com/golang/glog` — Leveled logging for the binary's own diagnostics (severity/verbosity flags, not structured)
- `github.com/prometheus/client_golang` — Prometheus metrics
