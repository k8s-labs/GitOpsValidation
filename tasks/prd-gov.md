# Product Requirements Document (PRD): GitOps Validation (gov)

## 1. Introduction/Overview

gov is a Kubernetes-native tool designed to validate that the actual state of a Kubernetes cluster matches the desired state as defined in a Flux GitOps repository. It ensures all Kubernetes resources specified in the GitOps repo are deployed and running as intended, providing continuous validation and reporting discrepancies. The tool is intended to be deployed as a pod within the cluster it validates, following Kubernetes and Go best practices.

## 2. Goals

- Ensure all Kubernetes resources defined in the Flux GitOps repository are present and healthy in the target cluster.
- Provide clear, actionable logs and Prometheus metrics for validation results.
- Allow skipping of resources not managed by the GitOps repo (e.g., system components).
- Operate as a headless service with no UI, suitable for automation and monitoring.

## 3. User Stories

- As a DevOps engineer, I want to be alerted if any resource defined in my GitOps repo is missing or unhealthy in the cluster, so I can take corrective action.
- As a platform team member, I want to see Prometheus metrics about validation status, so I can monitor cluster compliance over time.
- As an SRE, I want to exclude system components from validation, so that only managed resources are checked.

## 4. Functional Requirements

1. The system must accept configuration via command-line flags and environment variables, with command-line taking precedence.
2. The system must require the GitOps repo URL as a parameter (no default).
3. The system must support optional parameters: repo user (default: "gitops"), PAT (optional), branch (default: "main"), path (default: "./"), and wait time in seconds (default: 60).
4. The system must log a structured JSON message on startup.
5. The system must validate all parameters and exit with code 1 on error, logging the error.
6. The system must clone or verify the repo, checking out the specified branch and path, and pulling the latest changes.
7. The system must parse all Kubernetes resources defined in the specified path of the repo.
8. The system must validate that each resource is present and healthy in the cluster, using the Kubernetes API.
9. The system must validate all resource types found in the GitOps repo. If a resource type cannot be validated, gov must print a warning including the namespace, name, type, and other relevant details for further debugging.
10. For resources that define a healthcheck, gov should verify that the healthcheck was successful as part of validation.
11. The system must log information for each valid resource, warnings for extra resources in the cluster, and errors for missing/unhealthy resources.
12. The system must allow configuration to skip validation of resources not present in the GitOps repo (e.g., system components).
13. If the GitOps repo is unavailable, gov should log an error message and emit a Prometheus metric, then retry at the next configured sleep interval.
14. The system must expose Prometheus metrics for validation results (e.g., number of resources validated, number of discrepancies, last validation time).
15. The system must provide a `/healthz` HTTP endpoint for Kubernetes liveness/readiness monitoring.
16. The system must repeat validation at the configured interval until stopped.
17. The system must use Kubernetes best practices for startup, health checks, shutdown, and logging.

## 5. Non-Goals (Out of Scope)

- No support for GitOps tools other than Flux in the initial version.
- No web or CLI UI; all output is via logs and Prometheus metrics.
- No remediation or auto-healing of resources—gov is strictly read-only/validation.
- No validation of resources not defined in the GitOps repo, unless explicitly configured.

## 6. Design Considerations

- Structured JSON logging for all log levels (Info, Warning, Error, Fatal).
- Prometheus metrics endpoint for integration with monitoring systems.
- Designed for deployment via Kustomize as a Kubernetes deployment/pod.
- Should be lightweight and secure, following Go and Kubernetes best practices.

## 7. Technical Considerations

- Written in Go, using standard libraries and idioms.
- Uses the Kubernetes API to query cluster state.
- Only supports Flux GitOps repositories.
- Should handle intermittent network or API errors gracefully, with retries and clear error reporting.
- Should be able to run with minimal permissions required for read-only validation.

## 8. Success Metrics

- 100% of resources defined in the GitOps repo are present and healthy in the cluster.
- Prometheus metrics accurately reflect validation status and are scrapeable by monitoring systems.
- Clear, actionable logs for all validation runs.
- Minimal false positives/negatives in validation results.

## 9. Open Questions

1. What is the expected behavior if the cluster API is unavailable for an extended period?
