# GitOps Validation (gov)

This document contains the requirements to generate a PRD for the gov application written in Go

## Purpose

- gov is written in Go using standard go best practices and security
- a go.mod file is created and used
- use a local go reference instead of a fully qualified github go reference
- gov is usually deployed on the Kubernetes cluster that it is validating as a K8s deployment / pod
  - gov can also be run as a command line tool
- gov is deployed by Kustomize - initial version/overlay is 0.0.1
- gov should use K8s best practices for startup, healtchecks, shutdown, logging, etc
  - a K8s best practice is for main() to return a non-zero error code
- gov should use "zap" for json structured logging for Information, Warning, Error, and Fatal
- when gov runs from the command line, it assumes kubectl is setup correctly and has the necessary permissions to run against the default cluster
- from the command line, gov usually validates one time and exits
- when ran in K8s, it runs in daemon mode and sleeps between each validation

## API Endpoints

- API endpoints are only started in daemon mode
- must include Prometheus metrics
- must have a /healthz web endpoint for K8s to use
  - return 200 "pass" on success
  - return error code on failure
- must include /version
  - return 200 "0.0.1" on success (no newline)

## Parameters

- use the cobra library for command line parsing
- gov parameters can be environment variables or command line options - command line takes precedence, followed by environment variable, followed by default value
- Environment variables are of the form GOV_NAMESPACE
- Command line params are of the form --namespace or -n
- namespace (-n) - the Kubernetes namespace Flux is deployed to (default: flux-system)
- source (-s) - The Flux "source" repo (default: gitops)
- kustomization (-k) - the base Kustomization (default: gitops) - this listens to the flux-system/listeners directory in the proper cluster directory - for example ./clusters/tx-austin/flux-system/listeners
- sleep (-l) time (in seconds) - time to sleep between validations - default 60 - only used in daemon mode
- damon (-d) - run as a daemon - run the validations, then sleep for x seconds
- version (-v) - print the current version (without a new line) - i.e. 0.0.1

## Application Flow

- func main() returns an int - non-zero on error
- use json structured logging to log messages, warnings, errors, and fatal
- create a placeholder validation function that logs a validation success message
  - in --daemon mode, run validation, sleep for -l, loop
- gov starts and logs a start message
  - gov validates the parameters
    - gov logs any errors and returns 1
- log a stop message with error code or 0

## Docker

- use a multi-stage dockerfile that uses best practices for docker and security
  - docker cannot run as root
- use the latest Debian images
- the final image will need to have git and kustomize installed
- update the image using apt-get update and upgrade

## Kubernetes

- The kubernetes manifests for deploying gov should be created using Kustomize
  - the base directory should contain the namespace, service, and deployment
  - the overlays directory should contain "0.0.1" which is the current version
  - the gov structure should be in the clusters/tx-austin directory

- Add details on notification/alerting requirements.
    - the alerting is handled via processing the structured json logs using something like fluent bit forwarding to Azure Log Analytics
- Specify any compliance or regulatory requirements.
    - none currently
- Clarify if custom validation logic or plugins are needed.
    - not currently
- List all Kubernetes resource types to be validated.
    - namespace, service, deployment, pod, configmap, secrets, crds
- Define error handling and retry strategies.
    - gov should log the error using json structured logging and sleep until next iteration
- State requirements for metrics, observability, or integration with monitoring tools.
    - gov should expose standard Prometheus endpoint and metrics
