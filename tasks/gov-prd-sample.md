# Product Requirements Document (PRD): GitOps Validation (gov)

## Overview

gov is a tool designed to validate that a set of Flux GitOps manifests are deployed correctly on a Kubernetes cluster. It is typically deployed as a pod within the cluster it is validating. The application is written in Go, follows Go and Kubernetes best practices, and is deployed using Kustomize.

## Purpose
- Ensure that all resources defined in a GitOps repository are correctly deployed and running on the target Kubernetes cluster.
- Provide clear, structured logging for operational visibility and troubleshooting.
- Support flexible configuration via command line arguments and environment variables.

## Functional Requirements
1. Accept configuration via command line arguments and environment variables (command line > env var > default).
2. Validate all required parameters at startup; log errors and exit with code 1 if invalid.
3. Clone the specified repo using provided credentials; verify repo if directory exists.
4. Checkout the specified branch and pull the latest changes.
5. Change to the specified path within the repo; log and exit(1) on error.
6. Parse the GitOps manifests for namespaces, services, deployments, and pods.
7. Validate that each resource is deployed and running as specified on the Kubernetes cluster.
8. Log validation results using JSON structured logging (Information, Warning, Error, Fatal).
9. Support periodic validation based on the configured wait time.

## User Stories
- As a platform engineer, I want to ensure that all resources defined in my GitOps repo are deployed and healthy on my cluster.
- As an operator, I want to see clear, structured logs for all validation steps and errors.
- As a developer, I want to configure gov using either environment variables or command line arguments.

## Non-Functional Requirements
- Written in Go, following Go best practices and security guidelines.
- Deployed using Kustomize.
- Uses JSON structured logging.
- Follows Kubernetes best practices for lifecycle management (startup, health checks, shutdown, logging).

## Out of Scope
- Direct modification of cluster resources.
- UI/dashboard for validation results (logging only).

## Acceptance Criteria
- [ ] All required parameters are validated at startup.
- [ ] The application exits with code 1 on any fatal error.
- [ ] All validation results are logged in JSON format.
- [ ] The tool can be configured via both environment variables and command line arguments.
- [ ] The tool can be deployed as a Kubernetes pod using Kustomize.