# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This repository defines relevant concepts and types for Cluster Lifecycle APIs used in MCE (Multicluster Engine) and ACM (Advanced Cluster Management). It's a Go library that provides Kubernetes Custom Resource Definitions (CRDs) and client libraries for managing clusters in a multi-cluster environment.

## Common Development Commands

### Building and Verification
- `make build` - Build the project
- `make update` - Update all generated code (deepcopy, swagger docs, codegen, and CRDs)
- `make verify` - Run all verification scripts
- `make verify-scripts` - Run verification for deepcopy, swagger docs, CRDs, and codegen

### Individual Update Commands
- `make update-scripts` - Update deepcopy, swagger docs, and codegen
- `hack/update-deepcopy.sh` - Update deepcopy generation
- `hack/update-swagger-docs.sh` - Update swagger documentation
- `hack/update-codegen.sh` - Update client code generation
- `hack/update-crds.sh` - Update CRD generation

### Individual Verification Commands
- `hack/verify-deepcopy.sh` - Verify deepcopy generation
- `hack/verify-swagger-docs.sh` - Verify swagger documentation
- `hack/verify-crds.sh` - Verify CRD generation
- `hack/verify-codegen.sh` - Verify client code generation

## Architecture and Code Structure

### API Definitions

The repository is organized around five main API types, each in its own versioned directory:

1. **ManagedClusterAction** (`action/v1beta1/`) - Defines actions (Create/Update/Delete) to be executed on managed clusters
2. **ManagedClusterView** (`view/v1beta1/`) - Defines views of specific resources on managed clusters
3. **ManagedClusterInfo** (`clusterinfo/v1beta1/`) - Contains cluster metadata and status information
4. **ManagedClusterImageRegistry** (`imageregistry/v1alpha1/`) - Manages image registry overrides for managed clusters
5. **KlusterletConfig** (`klusterletconfig/v1alpha1/`) - Holds klusterlet configuration including proxy settings, install modes, and feature gates

### Generated Code Structure

The codebase heavily relies on code generation:

- **Client Libraries** (`client/*/`) - Generated Kubernetes clients, informers, and listers for each API
- **Deepcopy Methods** (`**/zz_generated.deepcopy.go`) - Runtime object deep copy implementations
- **Swagger Documentation** (`**/zz_generated.swagger_doc_generated.go`) - API documentation
- **CRD Manifests** - Generated from kubebuilder tags in type definitions

### Key Components

- **Constants** (`constants/`) - Shared constants and auto-import utilities
- **Helpers** (`helpers/`) - Utility functions for working with the APIs
- **Dependency Magnet** (`dependencymagnet/`) - Ensures proper vendoring of dependencies
- **Vendor** (`vendor/`) - Vendored dependencies
- **Tools** (`tools/`) - Development tooling dependencies

### Code Generation Flow

1. API types are defined in `*/v1beta1/types.go` or `*/v1alpha1/types.go` files
2. Kubebuilder annotations (e.g., `+genclient`, `+k8s:deepcopy-gen`) provide generation instructions
3. Controller-gen and other tools generate:
   - Deepcopy methods
   - Client libraries (clientsets, informers, listers)
   - CRD YAML manifests
   - Swagger documentation

### Important Patterns

- All API types follow standard Kubernetes conventions with TypeMeta, ObjectMeta, Spec, and Status
- Status conditions use standard metav1.Condition format
- Many fields are marked as deprecated with migration paths to newer alternatives
- Client libraries are generated separately for each API group
- The project uses OpenShift's build-machinery-go for common build patterns

## Development Workflow

1. **Making API Changes**: Modify types in the appropriate version directory
2. **Regenerate Code**: Run `make update` to regenerate all derived code
3. **Verify Changes**: Run `make verify` to ensure all generated code is consistent
4. **Testing**: No test framework is currently configured in this repository

## Dependencies and Tooling

- **Go Version**: 1.24.0
- **Controller-gen**: v0.17.3 (for CRD and code generation)
- **Runtime**: podman (configurable, can also use docker)
- **Build System**: Uses OpenShift's build-machinery-go for standardized builds
- **Key Dependencies**: Kubernetes client-go, controller-runtime, open-cluster-management.io/api

## Notes

- This is a library repository providing API definitions and generated clients
- Most files under `client/` and files named `zz_generated.*` are auto-generated
- Always run `make update` after modifying API type definitions
- The repository includes extensive vendor dependencies for offline builds