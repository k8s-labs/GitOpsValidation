# Product Requirements Document (PRD): GitOps Validation (gov)

## Overview

gov is a tool designed to validate that a set of Flux GitOps manifests are deployed correctly in a Kubernetes cluster. It is typically deployed as a pod within the cluster it monitors and ensures that the desired state described in GitOps manifests is accurately reflected in the cluster. The tool is written in Go, follows Kubernetes and Go best practices, and provides structured logging and health endpoints for observability.

## Purpose
- Solve the problem of ensuring Flux GitOps manifests are correctly applied in Kubernetes clusters.
- Users: Platform engineers, SREs, and DevOps teams responsible for cluster state and GitOps workflows.
- Value: Automated, repeatable validation of GitOps deployments, improved reliability, and observability.

## Functional Requirements
1. Validate that the specified Kubernetes namespace exists.
2. Retrieve and validate the Flux Source and Kustomization resources in the namespace.
3. Clone the specified GitOps repository using credentials from Kubernetes secrets.
4. Validate that the Flux Kustomization ran successfully and the manifests are applied.
5. Expose a /healthz endpoint for Kubernetes health checks (returns 200 "pass" or error code).
6. Expose a /validate endpoint to trigger immediate validation and return logs/results.
7. Expose Prometheus metrics endpoint for monitoring.
8. Support configuration via command line flags and environment variables (with precedence: CLI > env > default).
9. Log all actions, errors, and warnings using JSON structured logging.
10. Support periodic validation loop with configurable wait time.

## API Endpoints
- `/healthz` - Returns health status for Kubernetes probes.
- `/validate` - Triggers validation and returns results/logs.
- Prometheus metrics endpoint - Exposes standard metrics for monitoring.

## User Stories
- As a platform engineer, I want to ensure my GitOps manifests are correctly applied, so that my cluster state matches my desired configuration.
- As an SRE, I want to receive clear validation results and logs, so that I can quickly diagnose deployment issues.
- As a DevOps engineer, I want to trigger validation on demand, so that I can verify changes immediately after deployment.

## Non-Functional Requirements
- Use Go and Kubernetes best practices for startup, health checks, shutdown, and logging.
- Use JSON structured logging for all log levels.
- Secure handling of credentials (GitHub userId and PAT) via Kubernetes secrets.
- Performance: Validation should complete within 60 seconds for typical clusters.
- Scalability: Should support clusters with hundreds of resources.
- Observability: Expose Prometheus metrics and structured logs for integration with monitoring tools.

## Assumptions / Out of Scope
- Support for additional Git providers is not planned at this time
    - GitHub support requires https
    - For Private repos, it requires UserId and PAT
- Custom validation logic for specific resource types is not planned at this time
- Compliance or regulatory requirements (none currently)

## Acceptance Criteria
- [ ] The tool validates namespace, source, and kustomization existence and status.
- [ ] The /healthz endpoint returns correct status codes.
- [ ] The /validate endpoint triggers validation and returns logs.
- [ ] Prometheus metrics are exposed and scrapeable.
- [ ] All logs are JSON structured and include error, warning, info, and fatal levels.
- [ ] Configuration precedence is CLI > env > default.
- [ ] The tool handles errors gracefully and retries on next loop.

---

## Design Considerations
- Deployed as a Kubernetes pod using Kustomize.
- Follows Go and Kubernetes best practices.
- Logs are designed for consumption by log forwarders (e.g., fluent bit to Azure Log Analytics).

## Technical Considerations
- Uses Kubernetes API to retrieve resources and secrets.
- Clones GitHub repositories using credentials from secrets.
- Periodic validation loop with configurable interval.

## Success Metrics
- Reduction in configuration drift incidents.
- Decrease in support tickets related to GitOps deployment issues.
- High availability and reliability of validation results.

## Open Questions
- Are there any specific edge cases or error conditions not covered?
