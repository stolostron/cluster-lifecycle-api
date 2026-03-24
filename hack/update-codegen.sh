#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../../../k8s.io/code-generator)}

verify="${VERIFY:-}"

source "${CODEGEN_PKG}/kube_codegen.sh"

for group in action view clusterinfo imageregistry klusterletconfig clusterview; do
  kube::codegen::gen_client \
    --output-pkg "github.com/stolostron/cluster-lifecycle-api/client/${group}" \
    --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.txt" \
    --output-dir "${SCRIPT_ROOT}/client/${group}" \
    --one-input-api ${group} \
    --with-watch \
    .
done

# Generate OpenAPIModelName() accessor functions for all API packages annotated with
# +k8s:openapi-model-package. For non-clusterview packages, openapi-gen writes
# zz_generated.model_name.go into each package's own directory; we send the combined
# openapi definitions file to a temp location that is discarded.
OPENAPI_TMPDIR="$(mktemp -d)"
trap 'rm -rf "${OPENAPI_TMPDIR}"' EXIT

for pkg in \
  "action/v1beta1" \
  "view/v1beta1" \
  "clusterinfo/v1beta1" \
  "imageregistry/v1alpha1" \
  "klusterletconfig/v1alpha1" \
  "clusterview/v1alpha1"; do
  rm -f "${SCRIPT_ROOT}/${pkg}/zz_generated.model_name.go"
  # Call openapi-gen directly with only our package as input so that
  # --output-model-name-file writes exclusively to our package's directory.
  # Using kube::codegen::gen_openapi would silently add k8s.io/apimachinery
  # extra packages, causing openapi-gen to overwrite their vendor model_name
  # files (which have proper upstream license headers) with our empty boilerplate.
  go run -mod=vendor k8s.io/kube-openapi/cmd/openapi-gen \
    --output-file zz_generated.openapi.go \
    --output-model-name-file zz_generated.model_name.go \
    --go-header-file "${SCRIPT_ROOT}/hack/boilerplate.txt" \
    --output-dir "${OPENAPI_TMPDIR}" \
    --output-pkg "discard" \
    --report-filename /dev/null \
    "github.com/stolostron/cluster-lifecycle-api/${pkg}"
done
