## Relevant Files

- `internal/config/config.go` - Handles configuration via CLI flags and environment variables. (created/modified for config parsing, object definition, and documentation)
- `docs/configuration.md` - Documentation for all configuration options. (created/modified for config documentation)
- `internal/validation/flux.go` - Logic for retrieving and extracting config values from Flux Source and Kustomization. (created/modified for config population)
- `internal/config/config_test.go` - Unit tests for configuration logic. (created/modified for config tests)
- `cmd/gov/main.go` - Main entry point for the gov application.
- `internal/validation/namespace.go` - Logic for validating Kubernetes namespaces.
- `internal/git/clone.go` - Handles cloning of GitOps repositories using credentials from Kubernetes secrets.
- `internal/validation/kustomization.go` - Validates Flux Kustomization and manifest application.
- `internal/api/healthz.go` - Implements the /healthz endpoint.
- `internal/api/validate.go` - Implements the /validate endpoint.
- `internal/metrics/metrics.go` - Exposes Prometheus metrics endpoint.
- `internal/logging/logging.go` - Implements JSON structured logging.
- `internal/loop/loop.go` - Implements the periodic validation loop.
- `test/` - Directory for unit and integration tests for all modules.

### Notes

- Unit tests should be placed alongside or within the `test/` directory, matching the structure of the codebase.
- Use Go's built-in testing framework for all tests.

## Tasks

- [x] 1.0 Support Configuration via CLI Flags and Environment Variables
    - [x] 1.1 Implement configuration parsing with precedence: CLI > env > default
    - [x] 1.2 The config object must include fields for:
        - repo (url)
        - userId
        - pat
        - branch
        - path
    - [x] 1.3 Populate these config fields using values retrieved from the Flux Source and Flux Kustomization in step 3
    - [x] 1.4 Document all configuration options
    - [ ] 1.5 Write tests for configuration logic
- [ ] 2.0 Validate Kubernetes Namespace Existence
    - [ ] 2.1 Use Kubernetes API to check if the specified namespace exists
    - [ ] 2.2 Log an error and exit if the namespace does not exist
    - [ ] 2.3 Write unit tests for namespace validation logic
- [ ] 3.0 Retrieve and Validate Flux Source and Kustomization
    - [ ] 3.1 Use Kubernetes API to retrieve Flux Source in the namespace
    - [ ] 3.2 Retrieve Flux Kustomization in the namespace
    - [x] 3.3 Extract repo (url), userId, pat, branch from Flux Source and path from Flux Kustomization
    - [ ] 3.4 Validate that both resources exist and are correctly configured
    - [ ] 3.5 Log errors and exit if resources are missing or misconfigured
    - [ ] 3.6 Write unit tests for Flux resource validation
- [ ] 4.0 Clone GitOps Repository Using Kubernetes Secrets
    - [ ] 4.1 Retrieve GitHub credentials (UserId, PAT) from Kubernetes secrets
    - [ ] 4.2 Clone the specified repository using the credentials
    - [ ] 4.3 Handle public and private repo scenarios
    - [ ] 4.4 Log errors and exit if cloning fails
    - [ ] 4.5 Write unit tests for repository cloning logic
- [ ] 5.0 Validate Flux Kustomization and Manifest Application
    - [ ] 5.1 Check that the Flux Kustomization ran successfully
    - [ ] 5.2 Validate that all manifests are applied as expected
    - [ ] 5.3 Log results and errors
    - [ ] 5.4 Write unit tests for kustomization and manifest validation
- [ ] 6.0 Implement /healthz Endpoint for Health Checks
    - [ ] 6.1 Create HTTP handler for /healthz
    - [ ] 6.2 Return 200 "pass" on success, error code on failure
    - [ ] 6.3 Integrate with Kubernetes liveness/readiness probes
    - [ ] 6.4 Write tests for /healthz endpoint
- [ ] 7.0 Implement /validate Endpoint for On-Demand Validation
    - [ ] 7.1 Create HTTP handler for /validate
    - [ ] 7.2 Trigger validation logic and return logs/results
    - [ ] 7.3 Write tests for /validate endpoint
- [ ] 8.0 Expose Prometheus Metrics Endpoint
    - [ ] 8.1 Integrate Prometheus metrics library
    - [ ] 8.2 Expose standard metrics endpoint
    - [ ] 8.3 Write tests for metrics endpoint
- [ ] 9.0 Implement JSON Structured Logging
    - [x] 9.1 Use a structured logging library (e.g., zap, logrus)
    - [x] 9.2 Ensure all logs are JSON formatted and include error, warning, info, and fatal levels
    - [ ] 9.3 Write tests for logging output
- [ ] 10.0 Implement Periodic Validation Loop
    - [ ] 10.1 Implement loop to run validation at configurable intervals
    - [ ] 10.2 Ensure proper startup, shutdown, and error handling
    - [ ] 10.3 Write tests for periodic loop behavior
