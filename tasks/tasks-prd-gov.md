## Relevant Files

- `cmd/main.go` - Entry point for the gov application, handles CLI and daemon logic.
- `internal/validation/validator.go` - Contains core validation logic for Kubernetes resources.
- `internal/logging/logger.go` - Implements structured JSON logging using zap.
- `internal/api/server.go` - Implements /healthz, /version, and Prometheus metrics endpoints.
- `Dockerfile` - Multi-stage Dockerfile for building and running gov securely.
- `kustomize/base/deployment.yaml` - Kustomize manifest for base deployment.
- `kustomize/base/service.yaml` - Kustomize manifest for service.
- `kustomize/base/namespace.yaml` - Kustomize manifest for namespace.
- `kustomize/overlays/0.0.1/kustomization.yaml` - Overlay for version 0.0.1 deployment.
- `go.mod` - Go module definition for dependency management.

### Notes

- Unit tests should be placed alongside the code files they are testing (e.g., `validator_test.go` next to `validator.go`).
- Use `go test ./...` to run all Go tests.

## Tasks


- [ ] 1.0 Set up Go project structure and dependencies
  - [ ] 1.1 Initialize Go module and tidy dependencies in `go.mod`
  - [ ] 1.2 Create project directory structure (`cmd/`, `internal/`, etc.)
  - [ ] 1.3 Add local Go module references (avoid fully qualified GitHub references)
  - [ ] 1.4 Set up basic main.go entry point
  - [ ] 1.5 Add initial unit test files for structure

- [ ] 2.0 Implement structured JSON logging with zap
  - [ ] 2.1 Add zap dependency and configure logger
  - [ ] 2.2 Replace standard logging with zap in all modules
  - [ ] 2.3 Ensure logs include Information, Warning, Error, and Fatal levels
  - [ ] 2.4 Add unit tests for logging (where applicable)

- [ ] 3.0 Implement CLI and daemon modes with Cobra
  - [ ] 3.1 Add Cobra dependency and set up CLI structure
  - [ ] 3.2 Implement command-line flags and environment variable parsing (with correct precedence)
  - [ ] 3.3 Implement daemon mode (looping validation with sleep interval)
  - [ ] 3.4 Implement version flag/command
  - [ ] 3.5 Add unit tests for CLI and daemon logic

- [ ] 4.0 Implement API endpoints: /healthz, /version, and Prometheus metrics
  - [ ] 4.1 Set up HTTP server for API endpoints (only in daemon mode)
  - [ ] 4.2 Implement /healthz endpoint (200 "pass" on success, error code on failure)
  - [ ] 4.3 Implement /version endpoint (returns version string)
  - [ ] 4.4 Integrate Prometheus metrics endpoint
    - stub this out for now - /metrics returns 200, text/plain, not implemented
  - [ ] 4.5 Add unit tests for API endpoints

- [ ] 5.0 Implement core validation logic for Kubernetes resources
  - stub out this functionality for now - create a validation function that logs an is valid log
  - [ ] 5.1 Define validation interface and types for resources (namespace, service, deployment, pod, configmap, secrets, CRDs)
  - [ ] 5.2 Implement validation functions for each resource type
  - [ ] 5.3 Integrate validation logic into main application flow
  - [ ] 5.4 Add unit tests for validation logic

- [ ] 6.0 Create Dockerfile and Kustomize manifests for deployment
  - [ ] 6.1 Write multi-stage Dockerfile (non-root, Debian, installs git & kustomize)
    - Dockerfile must set ENV GOV_DAEMON=true and ENV GOV_SLEEP=10
  - [ ] 6.2 Add healthcheck to Dockerfile
  - [ ] 6.3 Create Kustomize base manifests: deployment, service, namespace
  - [ ] 6.4 Create Kustomize overlay for version 0.0.1
  - [ ] 6.5 Add documentation for deployment and configuration
  - [ ] 6.6 Add tests for deployment (e.g., validate manifests with kubeval or similar)
