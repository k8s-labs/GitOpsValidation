# Product Requirements Document (PRD): GitOps Validation (gov)

## Overview

gov is a Kubernetes-native application that validates the deployment of Flux GitOps manifests on a target Kubernetes cluster. Its primary goal is to ensure that the actual state of the cluster matches the desired state as defined in a GitOps repository, providing continuous validation, structured logging, and operational metrics for platform teams.

## Purpose
- Solves the problem of configuration drift by continuously validating that the cluster matches the desired state in the GitOps repository.
- Users are platform engineers, operators, and SREs responsible for cluster reliability and compliance.
- Provides value by automating validation, reducing manual checks, and surfacing issues early through logs and metrics.

## Functional Requirements
1. Accept configuration via command line arguments and environment variables, with command line taking precedence.
2. Support parameters: namespace (-n), source (-s), kustomization (-k), and wait time (seconds).
3. Use the Kubernetes API to validate the namespace, retrieve Flux source and kustomization, and extract repo URL, userId, branch, and PAT.
4. Clone and pull the GitOps repo using the provided credentials; change to ./gitops directory.
5. Validate that the namespace, Flux source, and kustomization exist and that the kustomization ran without issues.
6. Log information for each valid deployment and warnings for resources present on the cluster but not in the repo.
7. Repeat the validation process at the configured interval until stopped.
8. Use JSON structured logging for all log levels.
9. Expose Prometheus metrics and provide a /healthz web endpoint for Kubernetes health checks.
10. Follow Kubernetes best practices for startup, health checks, shutdown, and logging.
11. When run from the command line, assume kubectl is configured and has necessary permissions.

## User Stories
- As a platform engineer, I want to automatically validate that my cluster matches the desired state in my GitOps repo, so that I can ensure compliance and reliability.
- As an operator, I want to be alerted to any resources that are present on the cluster but not in the repo, so that I can investigate potential configuration drift.
- As a developer, I want to configure gov using environment variables or command line options, so that I can easily integrate it into different environments.
- As an SRE, I want to monitor gov's health and metrics using Prometheus and Kubernetes health checks.

## Non-Functional Requirements
- Must be performant and able to validate large clusters within a reasonable time window.
- Must use secure practices for handling credentials and secrets.
- Must be scalable to support multiple clusters and namespaces.
- Must be reliable and handle network/API errors gracefully.

## Out of Scope
- Providing a web UI or dashboard.
- Making changes to the cluster (gov is read-only/validation only).
- Supporting non-Flux GitOps tools or non-GitHub sources.
- Managing or rotating credentials.
- Alerting (handled via log forwarding systems such as FluentBit).
- Validating additional resource types beyond those specified.

## Acceptance Criteria
- [ ] All resources defined in the GitOps repo are validated as present and running on the cluster.
- [ ] All configuration errors are logged and result in a non-zero exit code.
- [ ] Prometheus metrics are available and accurate.
- [ ] /healthz endpoint is available and reflects application health.
- [ ] The system repeats validation at the configured interval until stopped.
