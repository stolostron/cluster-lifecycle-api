SHELL :=/bin/bash

all: build
.PHONY: all

# Ensure update-scripts are run before crd-gen so updates to Godoc are included in CRDs.
update-codegen-crds: update-scripts

# Include the library makefile
include $(addprefix ./vendor/github.com/openshift/build-machinery-go/make/, \
	golang.mk \
	targets/openshift/deps.mk \
	targets/openshift/crd-schema-gen.mk \
)

GO_PACKAGES :=$(addsuffix ...,$(addprefix ./,$(filter-out test/, $(filter-out vendor/,$(filter-out hack/,$(wildcard */))))))
GO_BUILD_PACKAGES :=$(GO_PACKAGES)
GO_BUILD_PACKAGES_EXPANDED :=$(GO_BUILD_PACKAGES)
# LDFLAGS are not needed for dummy builds (saving time on calling git commands)
GO_LD_FLAGS:=
CONTROLLER_GEN_VERSION :=v0.7.0

# $1 - target name
# $2 - apis
# $3 - manifests
# $4 - output
$(call add-crd-gen,actionv1beta1,./action/v1beta1,./action/v1beta1,./action/v1beta1)
$(call add-crd-gen,viewv1beta1,./view/v1beta1,./view/v1beta1,./view/v1beta1)
$(call add-crd-gen,clusterinfov1beta1,./clusterinfo/v1beta1,./clusterinfo/v1beta1,./clusterinfo/v1beta1)
$(call add-crd-gen,imageregistryv1beta1,./imageregistry/v1alpha1,./imageregistry/v1alpha1,./imageregistry/v1alpha1)
$(call add-crd-gen,inventoryv1beta1,./inventory/v1alpha1,./inventory/v1alpha1,./inventory/v1alpha1)

RUNTIME ?= podman
RUNTIME_IMAGE_NAME ?= openshift-api-generator

verify-scripts:
	bash -x hack/verify-deepcopy.sh
	bash -x hack/verify-swagger-docs.sh
	bash -x hack/verify-crds.sh
	bash -x hack/verify-codegen.sh
.PHONY: verify-scripts
verify: verify-scripts verify-codegen-crds

update-scripts:
	hack/update-deepcopy.sh
	hack/update-swagger-docs.sh
	hack/update-codegen.sh
.PHONY: update-scripts
update: update-scripts update-codegen-crds
