---
title: GitOps Validation (gov) PRD
author: bartr
date: 2025-06-13
---

# Product Requirements Document (PRD): GitOps Validation (gov)

## Overview

gov is a tool designed to validate that a set of Flux GitOps manifests are deployed correctly on a Kubernetes cluster. It is typically deployed as a pod within the cluster it is validating. The application is written in Go, follows Go and Kubernetes best practices, and is deployed using Kustomize.

## Purpose
- Ensure that all resources defined in a GitOps repository are correctly deployed and running on the target Kubernetes cluster.
- Provide clear, structured logging for operational visibility and troubleshooting.
- Support flexible configuration via command line arguments and environment variables.

## Functional Requirements
1. **Parameter Handling**
   - Accept configuration via command line arguments and environment variables (command line > env var > default).
   - Required parameters:
     - `--gov-repo` / `-r` or `GOV_REPO` (no default, required): GitOps repo URL (usually HTTPS).
   - Optional parameters with defaults:
     - User ID for repo: default `gitops`.
     - Personal Access Token (PAT): no default, not required for public repos.
     - Branch: default `main`.
     - Path within repo: default `./`.
     - Wait time (seconds): default `60`.

2. **Startup and Validation**
   - Log a structured startup message.
   - Validate all required parameters; log errors and exit with code 1 if invalid.

3. **Repository Operations**
   - Clone the specified repo using provided credentials.
   - If the repo directory exists, verify it matches the configured repo; otherwise, log error and exit(1).
   - Checkout the specified branch; log and exit(1) on error.
   - Pull the latest changes from the branch.
   - Change to the specified path within the repo; log and exit(1) on error.

4. **Validation Logic**
   - Parse the GitOps manifests for namespaces, services, deployments, and pods.
   - Validate that each resource is deployed and running as specified on the Kubernetes cluster.
   - Log validation results using JSON structured logging (Information, Warning, Error, Fatal).

5. **Operational Best Practices**
   - Use Kubernetes best practices for startup, health checks, shutdown, and logging.
   - Support periodic validation based on the configured wait time.

## User Stories
- As a platform engineer, I want to ensure that all resources defined in my GitOps repo are deployed and healthy on my cluster.
- As an operator, I want to see clear, structured logs for all validation steps and errors.
- As a developer, I want to configure gov using either environment variables or command line arguments.

## Non-Functional Requirements
- Written in Go, following Go best practices and security guidelines.
- Deployed using Kustomize.
- Uses JSON structured logging.
- Follows Kubernetes best practices for lifecycle management.

## Out of Scope
- Direct modification of cluster resources.
- UI/dashboard for validation results (logging only).

## Acceptance Criteria
- All required parameters are validated at startup.
- The application exits with code 1 on any fatal error.
- All validation results are logged in JSON format.
- The tool can be configured via both environment variables and command line arguments.
- The tool can be deployed as a Kubernetes pod using Kustomize.
