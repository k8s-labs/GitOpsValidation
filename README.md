

# gov

`gov` is a GitOps validation tool written in Go. It can be run as a CLI or as a Kubernetes deployment, providing validation, health checks, and metrics for GitOps-managed clusters.


## Table of Contents

- [Overview](#overview)
- [Configuration](#configuration)
- [Building and Running](#building-and-running)
- [API Endpoints](#api-endpoints)
- [Deployment](#deployment)
- [Logging and Observability](#logging-and-observability)
- [Development and Contribution](#development-and-contribution)
- [Examples](#examples)


## Overview

`gov` is designed to validate the state of a Kubernetes cluster managed by Flux. It can run as a one-shot CLI tool or as a daemon in Kubernetes, periodically validating resources and exposing health and metrics endpoints.


Key features:

- Runs as a CLI or Kubernetes deployment
- Validates resources (namespace, service, deployment, pod, configmap, secrets, CRDs)
- Structured JSON logging with zap
- Exposes `/healthz`, `/version`, and `/metrics` endpoints
- Configurable via CLI flags, environment variables, or defaults


## Configuration

Parameters can be set via CLI flags, environment variables, or will fall back to defaults. CLI flags take precedence over environment variables, which take precedence over defaults.

| Parameter      | CLI Flag         | Env Var                | Default        | Description                                      |
|---------------|------------------|------------------------|---------------|--------------------------------------------------|
| namespace     | --namespace, -n  | GOV_NAMESPACE          | flux-system    | Kubernetes namespace Flux is deployed to          |
| source        | --source, -s     | GOV_SOURCE             | gitops         | Flux source repo                                 |
| kustomization | --kustomization, -k | GOV_KUSTOMIZATION   | gitops         | Base Kustomization                               |
| sleep         | --sleep, -l      | GOV_SLEEP              | 60            | Sleep time in seconds between validations         |
| daemon        | --daemon, -d     | GOV_DAEMON             | false          | Run as daemon                                    |
| version       | --version, -v    | GOV_VERSION            | false          | Print version and exit                           |


### Precedence

1. CLI flag
2. Environment variable
3. Default value


## Building and Running

### Prerequisites

- Go 1.24+
- Docker (for container builds)
- Kubernetes cluster (for deployment)


### Build the CLI

```sh
go build -o gov ./cmd/gov
```


### Run as CLI (one-shot validation)

```sh
./gov --namespace my-ns --source my-repo
```

Or with environment variables:

```sh
export GOV_NAMESPACE=my-ns
export GOV_DAEMON=true
./gov
```


### Run in Daemon Mode

```sh
./gov --daemon --sleep 120
```


### Run Tests

```sh
go test ./...
```

## API Endpoints

When running in daemon mode (Kubernetes):

- `/healthz` - Returns `200 pass` if healthy
- `/version` - Returns the version string (e.g., `0.0.1`)
- `/metrics` - Exposes Prometheus metrics (stubbed in initial version)

## Deployment

### Kubernetes (Kustomize)

Kustomize manifests are in `clusters/tx-austin/gov/`.

To deploy the 0.0.1 overlay:

```sh
kubectl apply -k clusters/tx-austin/gov/overlays/0.0.1
```

Or from the gov directory:

```sh
kubectl apply -k .
# (if your KUBECONFIG context is correct and you are in clusters/tx-austin/gov)
```

### Docker

The Dockerfile builds a minimal image for Kubernetes deployment. See `.dockerignore` for excluded files.


## Logging and Observability

- Uses zap for structured JSON logs (info, warning, error, fatal)
- Logs startup, shutdown, validation results, and errors
- Metrics endpoint for Prometheus scraping


## Development and Contribution

- Code is organized in `cmd/` (entry point) and `internal/` (modules)
- Use Go modules (`go.mod`, `go.sum`)
- Run `go test ./...` for all tests
- Add new features with tests and documentation
- Follow Go best practices and idiomatic style

### Directory Structure

- `cmd/gov/main.go` - CLI entry point
- `internal/config/` - Configuration parsing
- `internal/logging/` - Logging setup
- `internal/validation/` - Validation logic
- `internal/server/` - HTTP endpoints
- `clusters/` - Kustomize manifests for Kubernetes


## Examples

### CLI Example

```sh
./gov --namespace flux-system --source gitops --kustomization gitops
```

### Kubernetes Example

```sh
kubectl apply -k clusters/tx-austin/gov/overlays/0.0.1
```
