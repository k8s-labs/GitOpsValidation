# GitOps Validation (gov)

This document contains the requirements to generate a PRD for the GitOps Validation (gov) application

## Purpose

- gov validates that a set of Flux GitOps manifests are deployed correctly
- gov is usually deployed on the Kubernetes cluster that it is validating as a K8s deployment / pod
- gov is written in Go using standard go best practices and security
- gov is deployed by Kustomize
- gov should use K8s best practices for startup, healtchecks, shutdown, logging, etc
- gov should use json structured logging for Information, Warning, Error, and Fatal
- gov should include Prometheus metrics
- gov should have a /healthz web endpoint for K8s to use
- gov usually runs as a pod / deployment in the K8s cluster it is monitoring
- when gov runs from the command line, it assumes kubectl is setup correctly and has the necessary permissions to run against the default cluster

## Parameters

- gov parameters can be environment variables or command line options - command line takes precedence, followed by environment variable, followed by default value
- Environment variables are of the form GOV_NAMESPACE
- Command line params are of the form --namespace or -n
- namespace (-n) - the Kubernetes namespace Flux is deployed to (default: flux-system)
- source (-s) - The Flux "source" repo (default: gitops) - must be GitHub via https - can be a public repo (no PAT required) or private repo (userid and pat required)
- kustomization (-k) - the base Kustomization (default: flux-listeners) - this listens to the flux-system/listeners directory in the proper cluster directory - for example ./clusters/tx-austin/flux-system/listeners
- wait time (in seconds) - time to wait between validations - default 60

## Application Flow

- gov starts and logs a starting message
  - gov validates the parameters
    - gov logs any error messages and exits with a value of 1
- gov uses the k8s API to validate the namespace exists
  - use the k8s API to get the Flux Source
    - save the repo URL, userId, branch, and PAT (will be in a k8s secret) into a struct for later use
  - use the k8s API to get the Flux Kustomization
    - save the path to a struct for later use
  - log any error and exit(1)
- clone and pull the GitOps repo
  - use json structured logging to log messages, warnings, errors, and fatal
  - if the gitops directory doesn't exist gov clones the repo to the gitops directory using the repo, branch, user, and PAT values
  - change the current directory to ./gitops
  - log any error and exit(1)
- validation loop
  - use json structured logging to log messages, warnings, errors, and fatal
  - ensure ./gitops is the current directory
  - git pull to get the latest from the repo
  - use the k8s API to
    - validate the namespace exists in the K8s cluster
    - validate the flux source exists in the namespace
    - validate the flux kustomization exists in the namespace
    - validate the flux kustomization ran without issues
  - sleep for wait time seconds and then repeats the validation loop until stopped

- gov only supports GitHub via https using userId and optional PAT
- No additional resource types are included in the validation?
- Alerting is handled via FluentBit or similar log forwarding to a log system such as Azure Log Analytics
