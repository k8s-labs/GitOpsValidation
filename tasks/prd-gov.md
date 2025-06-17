# Product Requirements Document (PRD): gov

## Overview

gov is a GitOps Validation Framework (fx) designed to validate Kubernetes clusters. It can be deployed as a Kubernetes pod or run as a command line tool. The application is written in Go, follows best practices for security and logging, and is deployed using Kustomize. The initial version is 0.0.1.

## Purpose
- Solves the problem of validating Kubernetes clusters for correct configuration and deployment.
- Users: Platform engineers, SREs, DevOps teams managing Kubernetes clusters with GitOps (Flux).
- Value: Provides automated, repeatable validation of cluster state, improving reliability and compliance.

## Functional Requirements
1. Deployable as a Kubernetes pod/deployment or as a CLI tool.
2. Use Go best practices for security, startup, health checks, shutdown, and logging.
3. Use "zap" for JSON structured logging (Info, Warning, Error, Fatal).
4. Support command line and environment variable configuration (command line > env var > default).
5. Validate parameters at startup and log errors.
6. Provide a placeholder validation function that logs success.
7. In daemon mode, run validation, sleep, and repeat.
8. Log start and stop messages with error codes.

## API Endpoints
- `/metrics`: Exposes Prometheus metrics for monitoring.
- `/healthz`: Returns 200 "pass" on success, error code on failure (for K8s health checks).
- `/version`: Returns 200 with version string (e.g., "0.0.1").

## User Stories
- As a platform engineer, I want to validate my cluster's configuration automatically, so that I can ensure compliance and reliability.
- As an SRE, I need health and metrics endpoints, so that I can monitor the application's status and performance.
- As a DevOps engineer, I want to run validations on demand or in daemon mode, so that I can integrate checks into CI/CD workflows.

## Non-Functional Requirements
- Performance: Should validate clusters quickly and efficiently.
- Security: Follows Go and Kubernetes security best practices.
- Scalability: Can be run in multiple clusters or as a CLI for different environments.
- Logging: Uses structured JSON logging for all events.

## Out of Scope
- Implementing actual validation logic (only a placeholder is required initially).
- Support for non-Kubernetes environments.
- UI/dashboard features.

## Acceptance Criteria
- [ ] Application can be deployed as a pod or run as a CLI.
- [ ] All configuration parameters work as described (CLI, env vars, defaults).
- [ ] Health, metrics, and version endpoints are available and function as specified.
- [ ] Logging is structured and includes all required levels.
- [ ] Placeholder validation function logs success.
- [ ] Daemon mode works as described.
