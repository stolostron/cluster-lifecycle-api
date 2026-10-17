# Architecture

## Purpose

`cluster-lifecycle-api` is a Go library of Kubernetes API definitions and
supporting client helpers for Red Hat Advanced Cluster Management and
Multicluster Engine. It does not contain a controller or a standalone server.
Consumers import its typed APIs, generated clients, CRDs, and helper packages.

## API packages

Each versioned package owns source API types, group registration, and generated
artifacts:

- `action/v1beta1`: describes actions executed against resources on managed clusters.
- `view/v1beta1`: describes reads of selected resources on managed clusters.
- `clusterinfo/v1beta1`: represents hub-observed managed-cluster information,
  including node, distribution, version, endpoint, and condition data.
- `imageregistry/v1alpha1`: configures image registry replacement for selected
  managed clusters through Placement references.
- `klusterletconfig/v1alpha1`: configures klusterlet installation, hub API
  connectivity, proxy and trust settings, placement, feature gates, and related
  agent behavior.

The packages use Kubernetes API conventions and annotations. Source structs and
markers are the authoritative inputs; deepcopy code, Swagger documentation,
model-name files, generated clients, and CRDs are derived outputs.

## Helper layer

`helpers/` contains reusable behavior built on the public API types:

- `helpers/imageregistry` reads managed-cluster annotations and Kubernetes
  Secrets to resolve image pull secrets and registry overrides.
- `helpers/tlsprofile` watches and resolves TLS profile configuration.
- `helpers/klusterletconfig` merges klusterlet configuration values.
- `helpers/localcluster` and `helpers/userpermission` provide focused utilities
  for consumers.

These helpers are library code. They receive clients or data from their callers
and do not own reconciliation loops or process lifecycle.

## Generation and verification flow

API source changes flow through the repository's generation scripts:

1. API structs and Kubernetes/controller-gen markers are edited in a versioned package.
2. `make update` regenerates deepcopy, Swagger, codegen, and related CRD output.
3. `make verify` reruns checks in `hack/verify-*.sh` and fails when generated
   output differs from the checked-in result.
4. Downstream ACM/MCE repositories consume the module and its generated API
   artifacts.

The Makefile includes OpenShift build machinery and uses the vendored dependency
tree. This keeps generation and verification consistent with the repository's CI
workflow.

## Compatibility boundaries

The versioned API packages and serialized JSON fields are compatibility
boundaries for downstream Kubernetes clients. Changes to types, JSON tags,
annotations, markers, or generated CRDs should be treated as API changes and
validated with the complete generation and verification flow.
