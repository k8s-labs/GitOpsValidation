# Product Requirements Document: GitOps Validation (gov)

## Introduction/Overview

The GitOps Validation (gov) application is a monitoring solution designed to validate that Flux GitOps manifests are deployed correctly within a Kubernetes cluster. It serves as a continuous validation tool that ensures the integrity and proper functioning of GitOps deployments by regularly checking the state of Flux resources against their expected configurations. This tool addresses the common challenge of detecting drift or failed deployments in GitOps-managed environments.

## Goals

1. Provide real-time validation of Flux GitOps manifests deployed in a Kubernetes cluster
2. Detect and report any discrepancies between expected and actual states of GitOps resources
3. Offer robust monitoring capabilities with structured logging and metrics exposure
4. Operate with minimal resource overhead within the target Kubernetes cluster
5. Support standard Kubernetes operational practices (health checks, graceful shutdown, etc.)
6. Allow flexible configuration through both environment variables and command-line parameters

## User Stories

- As a DevOps engineer, I want to automatically validate that my Flux GitOps configurations are deployed correctly, so that I can quickly identify and resolve any issues in my GitOps workflow.
- As a platform administrator, I need to know when Flux deployments have failed or drifted from their expected state, so that I can maintain system reliability.
- As a Kubernetes operator, I want to monitor the health of my GitOps pipeline, so that I can ensure continuous delivery works as expected.
- As a developer, I want insight into the state of GitOps deployments, so that I can understand if my changes were successfully applied to the cluster.

## Functional Requirements

1. The system must provide a Go application that can be deployed as a Kubernetes pod within the cluster it is validating.
2. The system must validate that the specified Flux namespace exists in the Kubernetes cluster.
3. The system must validate that the specified Flux Source resource exists and is properly configured.
4. The system must validate that the specified Flux Kustomization resource exists and ran without issues.
5. The system must support connecting to GitHub repositories via HTTPS, with support for both public and private repositories (with PAT authentication).
6. The system must clone and pull from the GitOps repository to get the latest configuration.
7. The system must perform validation checks at a configurable interval (default: 60 seconds).
8. The system must provide JSON structured logging for Information, Warning, Error, and Fatal events.
9. The system must expose Prometheus metrics for monitoring purposes.
10. The system must provide a `/healthz` web endpoint for Kubernetes health checks.
11. The system must support configuration through environment variables (prefixed with `GOV_`) and command-line options.
12. The system must recognize command-line parameters taking precedence over environment variables, which take precedence over default values.
13. The system must support the following parameters:
    - Namespace (`-n`, `--namespace`, `GOV_NAMESPACE`): The Kubernetes namespace Flux is deployed to (default: `flux-system`)
    - Source (`-s`, `--source`, `GOV_SOURCE`): The Flux "source" repo (default: `gitops`)
    - Kustomization (`-k`, `--kustomization`, `GOV_KUSTOMIZATION`): The base Kustomization (default: `flux-listeners`)
    - Wait time (`--wait`, `GOV_WAIT`): Time in seconds to wait between validations (default: 60)
14. The system must fail gracefully with appropriate error messages and non-zero exit code when validation checks fail.
15. The system must be deployable using Kustomize.

## Non-Goals (Out of Scope)

1. The system will not support Git providers other than GitHub (no GitLab, Bitbucket, etc.).
2. The system will not support SSH-based Git authentication.
3. The system will not modify or apply Kubernetes resources; it is validation-only.
4. The system will not directly send alerts (alerting is handled via external systems like FluentBit).
5. The system will not validate resources beyond Flux Source and Kustomization objects.
6. The system will not provide a web UI for visualization of validation results.
7. The system will not manage or update the Flux installation itself.

## Design Considerations

- The application should follow Go best practices for code organization and project structure.
- The application should use a clean command-line interface with clear help documentation.
- Logs should be structured in JSON format to facilitate integration with log aggregation tools.
- Metrics should follow Prometheus best practices for naming and organization.

## Technical Considerations

1. The application should be written in Go using standard Go best practices and security considerations.
2. The application should use the official Kubernetes Go client library for interacting with the Kubernetes API.
3. The application should use a common/standard Go Git client library for repository operations.
4. The application should be containerized following best practices for minimal container images.
5. The application should handle Kubernetes credentials through standard mechanisms (service account, kubeconfig).
6. The application should have appropriate resource requests and limits defined in its deployment manifest.
7. The application should implement proper error handling and retry logic for transient failures.

## Success Metrics

1. Zero undetected GitOps deployment failures in monitored clusters.
2. Mean time to detect (MTTD) GitOps issues under 2 minutes.
3. Resource utilization under 100MB memory and 0.1 CPU cores during normal operation.
4. Successful operation in production environments for at least 30 days without restarts.
5. Complete coverage of all specified validation checks.

## Open Questions

1. Should the tool support additional Git providers beyond GitHub in future versions?
2. Should the tool provide webhook capabilities for immediate notification of validation failures?
3. What specific Prometheus metrics should be exposed for optimal monitoring?
4. Should the tool validate additional Flux resource types in future versions?
5. Is there a need for a dry-run mode to validate configurations without affecting the system?
