# Product Requirements Document (PRD): GitOps Validation (gov)

## Introduction/Overview

gov is a Kubernetes-native application that validates the deployment of Flux GitOps manifests on a target Kubernetes cluster. Its primary goal is to ensure that the actual state of the cluster matches the desired state as defined in a GitOps repository, providing continuous validation, structured logging, and operational metrics for platform teams.

## Goals

- Validate that all namespaces, services, deployments, and pods defined in the GitOps repository are present and running on the cluster.
- Detect and log any resources present on the cluster but not defined in the GitOps repository.
- Provide JSON structured logs for information, warnings, errors, and fatal events.
- Expose Prometheus metrics for monitoring.
- Provide a /healthz endpoint for Kubernetes health checks.
- Support flexible configuration via command line arguments and environment variables.
- Operate reliably as a Kubernetes deployment, following K8s best practices for lifecycle and security.

## User Stories

- As a platform engineer, I want to automatically validate that my cluster matches the desired state in my GitOps repo, so that I can ensure compliance and reliability.
- As an operator, I want to be alerted to any resources that are present on the cluster but not in the repo, so that I can investigate potential configuration drift.
- As a developer, I want to configure gov using environment variables or command line options, so that I can easily integrate it into different environments.
- As an SRE, I want to monitor gov's health and metrics using Prometheus and Kubernetes health checks.

## Functional Requirements

1. The system must accept configuration via command line arguments and environment variables, with command line taking precedence.
2. Environment variables must be of the form GOV_NAMESPACE, GOV_SOURCE, etc.; command line params must be of the form --namespace/-n, --source/-s, etc.
3. The system must support the following parameters:
   - namespace (-n): Kubernetes namespace Flux is deployed to (default: flux-system)
   - source (-s): Flux source repo (default: gitops, must be GitHub via https)
   - kustomization (-k): Base Kustomization (default: flux-listeners)
   - wait time: seconds between validations (default: 60)
4. On startup, gov must log a starting message and validate all parameters, exiting with error code 1 on failure.
5. The system must use the Kubernetes API to validate the namespace, retrieve the Flux source and kustomization, and extract repo URL, userId, branch, and PAT (from k8s secret).
6. The system must clone and pull the GitOps repo using the provided credentials. If the gitops directory doesn't exist, clone it; otherwise, pull the latest changes.
7. The system must change the current directory to ./gitops and log errors if unsuccessful.
8. The system must validate that the namespace, Flux source, and kustomization exist and that the kustomization ran without issues.
9. The system must log information for each valid deployment and warnings for resources present on the cluster but not in the repo.
10. The system must repeat the validation process at the configured interval until stopped.
11. The system must use JSON structured logging for all log levels.
12. The system must expose Prometheus metrics.
13. The system must provide a /healthz web endpoint for Kubernetes health checks.
14. The system must follow Kubernetes best practices for startup, health checks, shutdown, and logging.
15. When run from the command line, gov must assume kubectl is configured and has necessary permissions.

## Non-Goals (Out of Scope)

- Providing a web UI or dashboard.
- Making changes to the cluster (gov is read-only/validation only).
- Supporting non-Flux GitOps tools or non-GitHub sources.
- Managing or rotating credentials.
- Alerting (handled via log forwarding systems such as FluentBit).
- Validating additional resource types beyond those specified.

## Design Considerations

- Should be deployed as a Kubernetes deployment/pod using Kustomize.
- Should use standard Go and Kubernetes libraries for reliability and maintainability.
- Should be stateless and idempotent between validation cycles.
- Should be compatible with Prometheus and Kubernetes health checks.

## Technical Considerations

- Must integrate with Kubernetes API for resource validation.
- Must handle network or API errors gracefully, with retries and clear error logging.
- Must be compatible with standard Kubernetes RBAC and security practices.
- Must support both in-cluster and command-line operation modes.

## Success Metrics

- 100% of resources defined in the GitOps repo are validated as present and running on the cluster.
- All configuration errors are logged and result in a non-zero exit code.
- Prometheus metrics are available and accurate.
- /healthz endpoint is available and reflects application health.
- Reduction in configuration drift incidents reported by platform teams.

## Open Questions

- Should gov support additional Git providers (e.g., GitLab, Bitbucket) or only GitHub?
- Are there specific resource types (e.g., CRDs) that should be included in validation?
- Should there be integration with alerting systems (e.g., Slack, email) for warnings/errors, or is log forwarding sufficient?
- Should gov support custom validation plugins or rules in the future?
