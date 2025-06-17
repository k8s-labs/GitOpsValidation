# Product Requirements Document (PRD) – gov

## Overview

gov is a Go application designed to validate Kubernetes clusters using best practices for security, deployment, and observability. It can run as a command-line tool or as a daemonized service within a Kubernetes cluster, providing validation, health checks, and metrics.

## Purpose
- **Problem Solved:** Ensures Kubernetes clusters adhere to best practices and are properly configured, reducing operational risk.
- **Users:** Platform engineers, SREs, DevOps teams managing Kubernetes clusters.
- **Value:** Automated, repeatable validation of cluster state, with structured logging and metrics for monitoring and alerting.

## Functional Requirements
1. Must be written in Go, using standard Go best practices and security.
2. Use a local Go module reference (not a fully qualified GitHub reference).
3. Deployable as a Kubernetes deployment/pod (via Kustomize) or as a CLI tool.
4. Provide structured JSON logging using the "zap" library for all log levels.
5. Support both command-line and environment variable configuration (command-line takes precedence).
6. Validate Kubernetes resources: namespace, service, deployment, pod, configmap, secrets, CRDs.
7. Expose Prometheus metrics and a /healthz endpoint (200 "pass" on success, error code on failure).
8. Expose a /version endpoint (returns version string, e.g., "0.0.1").
9. Support daemon mode (looping validation with sleep interval).
10. Use Cobra for command-line parsing.
11. Provide Dockerfile (multi-stage, non-root, Debian base, installs git & kustomize).
12. Provide Kustomize manifests for deployment (base and overlays).

## API Endpoints
- **/healthz** – Returns 200 "pass" if healthy, error code otherwise (for K8s liveness/readiness probes).
- **/version** – Returns the current version string (e.g., "0.0.1").
- **Prometheus metrics** – Standard metrics endpoint for monitoring.

## User Stories
- As a platform engineer, I want to validate my Kubernetes cluster automatically, so that I can ensure compliance with best practices.
- As an SRE, I need structured logs and metrics, so that I can monitor and alert on cluster health.
- As a DevOps engineer, I want to run gov as a CLI or in-cluster daemon, so that I can use it flexibly in different environments.

## Non-Functional Requirements
- Must use structured JSON logging (zap).
- Must not run as root in Docker.
- Must use latest Debian images and update/upgrade packages.
- Must expose Prometheus metrics for observability.
- Must handle errors by logging and retrying after a sleep interval.
- Must follow Kubernetes best practices for startup, health checks, shutdown, and logging.

## Out of Scope
- Custom validation logic or plugins (not required initially).
- Compliance or regulatory requirements (none currently).
- Notification/alerting is handled externally via log forwarding (e.g., fluent bit to Azure Log Analytics).

## Acceptance Criteria
- [ ] Application is written in Go with local module references.
- [ ] CLI and daemon modes are supported with correct parameter precedence.
- [ ] Structured JSON logging is implemented using zap.
- [ ] /healthz and /version endpoints are available in daemon mode.
- [ ] Prometheus metrics endpoint is available.
- [ ] Dockerfile follows best practices (multi-stage, non-root, Debian, installs git & kustomize).
- [ ] Kustomize manifests are provided for deployment.
- [ ] Application validates all required Kubernetes resource types.
- [ ] Error handling and retry logic are implemented as specified.
