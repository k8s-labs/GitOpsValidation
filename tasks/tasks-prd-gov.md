## Relevant Files

- `cmd/gov/main.go` - Entry point for the gov application, handles startup, argument parsing, logs a structured startup message, and runs periodic validation.
- `cmd/gov/main_test.go` - Unit test for startup and validation logic (integration with config).
- `internal/config/config.go` - Defines the Config struct, implements command line and environment variable parsing with correct precedence, and adds validation logic for required parameters.
- `internal/git/repo.go` - Logic for cloning the GitOps repository, verifying the repo directory, branch checkout, pulling latest changes, and changing to the specified path with error handling.
- `internal/validator/validator.go` - Parses GitOps manifests, compares definitions to cluster state, and logs validation results.
- `internal/logger/logger.go` - JSON structured logging implementation, used for structured startup and error messages.
- `internal/k8s/client.go` - Kubernetes client interactions for resource validation and health checks.
- `deploy/kustomize/` - Kustomize manifests for deploying gov as a Kubernetes pod.
- `internal/config/config_test.go` - Unit tests for configuration handling, including environment variable and validation logic.
- `internal/git/repo_test.go` - Unit tests for repository operations, including path and repo verification logic.
- `internal/validator/validator_test.go` - Unit tests for manifest parsing and validation logic.
- `internal/logger/logger_test.go` - Unit tests for logging (Info, Warn, Error).
- `internal/k8s/client_test.go` - Unit tests for Kubernetes client health check logic.

### Notes

- Unit tests should typically be placed alongside the code files they are testing (e.g., `validator.go` and `validator_test.go` in the same directory).
- Use `npx jest [optional/path/to/test/file]` to run tests. Running without a path executes all tests found by the Jest configuration.

## Tasks

- [ ] 1.0 Implement Configuration Handling
  - [x] 1.1 Define configuration struct for all parameters (repo URL, user, PAT, branch, path, wait time)
  - [x] 1.2 Implement parsing of command line arguments using a Go CLI library (e.g., cobra or flag)
  - [x] 1.3 Implement environment variable parsing with correct precedence
  - [x] 1.4 Add validation logic for required parameters (e.g., repo URL)
  - [x] 1.5 Write unit tests for configuration handling

- [ ] 2.0 Implement Application Startup and Parameter Validation
  - [x] 2.1 Log a structured startup message using JSON logging
  - [x] 2.2 Validate all required parameters at startup
  - [x] 2.3 Exit with code 1 and log error if validation fails
  - [x] 2.4 Write unit tests for startup and validation logic

- [ ] 3.0 Implement Repository Operations
  - [x] 3.1 Implement logic to clone the GitOps repository using provided credentials
  - [x] 3.2 If repo directory exists, verify it matches the configured repo
  - [x] 3.3 Implement branch checkout and error handling
  - [x] 3.4 Implement git pull for latest changes
  - [x] 3.5 Change to the specified path within the repo and handle errors
  - [x] 3.6 Write unit tests for repository operations

- [ ] 4.0 Implement Validation Logic for Kubernetes Resources
  - [x] 4.1 Parse GitOps manifests for namespaces, services, deployments, and pods
  - [x] 4.2 Implement Kubernetes client logic to check resource status
  - [x] 4.3 Compare manifest definitions to actual cluster state
  - [x] 4.4 Log validation results (Information, Warning, Error, Fatal) in JSON format
  - [x] 4.5 Write unit tests for validation logic

- [ ] 5.0 Implement Logging and Operational Best Practices
  - [x] 5.1 Implement JSON structured logging for all log messages
  - [x] 5.2 Add Kubernetes best practices for startup, health checks, and shutdown
  - [x] 5.3 Implement periodic validation based on wait time parameter
  - [x] 5.4 Write unit tests for logging and operational logic
