## Relevant Files

- `cmd/gov/main.go` - Entry point for the gov CLI and Kubernetes deployment logic. (Created for sub-task 1.1)
- `internal/config/config.go` - Handles configuration parsing from CLI, environment variables, and defaults. (Created and implemented for task 2)
- `internal/logging/logging.go` - Sets up structured JSON logging using zap. (Created and implemented for task 3)
- `internal/validation/validation.go` - Contains the placeholder validation logic. (Created and implemented for task 4)
- `internal/server/server.go` - Implements the HTTP endpoints for /metrics, /healthz, and /version. (Created and implemented for task 5)
- `Dockerfile` - Containerizes the application for Kubernetes deployment.
- `README.md` - Documentation for usage, configuration, and deployment.
- `go.mod` - Go module definition file. (Created for sub-task 1.6)

### Notes

- Unit tests should be placed alongside the code files they are testing (e.g., `validation.go` and `validation_test.go` in the same directory).
- Use `go test ./...` to run all tests in the project.
- repo root is ~/gov


## Tasks

- [ ] 1.0 Implement application entry point and CLI logic
    - [ ] 1.1 Set up Go module and project structure
    - [ ] 1.2 Create `main.go` with basic CLI entry point - func main() must return non-zero int on error
    - [ ] 1.3 Integrate cobra for CLI command parsing
    - [ ] 1.4 Add version flag/command
    - [ ] 1.4.1 the web server should only start in daemon mode
    - [ ] 1.5 Handle exit codes and error propagation
    - [ ] 1.6 run `go build` until all errors are resolved

- [ ] 2.0 Implement configuration parsing (CLI, env vars, defaults) with cobra
    - [ ] 2.1 Define all configuration parameters (namespace, source, kustomization, sleep, daemon, version)
    - [ ] 2.2 Implement precedence: CLI > env var > default
    - [ ] 2.3 Document all parameters in README
    - [ ] 2.4 Add validation for required parameters

- [ ] 3.0 Implement structured JSON logging with zap
    - [ ] 3.1 Add zap as a dependency
    - [ ] 3.2 Set up logger for Info, Warning, Error, Fatal
    - [ ] 3.3 Integrate logging throughout the app (startup, shutdown, errors)
    - [ ] 3.4 Ensure logs are JSON-formatted

- [ ] 4.0 Implement placeholder validation logic
    - [ ] 4.1 Create validation function/module
    - [ ] 4.2 Log validation success/failure
    - [ ] 4.3 Integrate validation into main app flow

- [ ] 5.0 Implement HTTP endpoints: /metrics, /healthz, /version
    - [ ] 5.1 Set up HTTP server in a separate goroutine
    - [ ] 5.2 Implement /healthz endpoint (200 "pass" or error)
    - [ ] 5.3 Implement /version endpoint (returns version string)
    - [ ] 5.4 Integrate Prometheus metrics and expose /metrics

- [ ] 6.0 Implement daemon mode (looping validation with sleep)
    - [ ] 6.1 Add daemon flag/parameter
    - [ ] 6.2 Implement loop: run validation, sleep, repeat
    - [ ] 6.3 Ensure graceful shutdown and signal handling

- [ ] 7.0 Create Kustomize manifests and Dockerfile for deployment
    - [ ] 7.1 Write Dockerfile and .dockerignore for gov app
    - [ ] 7.2 Create Kustomize base and overlay manifests
    - [ ] 7.3 Add Kubernetes deployment, service, and health checks
    - [ ] 7.4 Document deployment steps

- [ ] 8.0 Write documentation and usage instructions
    - [ ] 8.1 Document CLI usage and parameters
    - [ ] 8.2 Document environment variable configuration
    - [ ] 8.3 Add examples for running in CLI and Kubernetes
    - [ ] 8.4 Document endpoints and expected responses
    - [ ] 8.5 Add contribution and development notes
