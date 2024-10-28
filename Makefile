# Set these to the desired values
ARTIFACT_ID=k8s-ces-control
VERSION=0.11.1
GOTAG=1.22.4
LINT_VERSION=v1.58.2
STAGE?=production
LOG_LEVEL?=info

# Setting SHELL to bash allows bash commands to be executed by recipes.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

## Image URL to use all building/pushing image targets
IMAGE?=cloudogu/${ARTIFACT_ID}:${VERSION}

MAKEFILES_VERSION=9.3.2
.DEFAULT_GOAL:=default
GENERATION_TARGET_DIR=generated
GENERATION_SOURCE_DIR=grpc-protobuf
INTEGRATION_TEST_NAME_PATTERN=.*_inttest$$
# make sure to create a statically linked binary
GO_BUILD_FLAGS=-mod=vendor -a -tags netgo,osusergo $(LDFLAGS) -o $(BINARY)

K8S_RESOURCE_DIR=${WORKDIR}/k8s

include build/make/variables.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/digital-signature.mk
include build/make/self-update.mk
include build/make/release.mk
include build/make/bats.mk

K8S_COMPONENT_SOURCE_VALUES = ${HELM_SOURCE_DIR}/values.yaml
K8S_COMPONENT_TARGET_VALUES = ${HELM_TARGET_DIR}/values.yaml
HELM_PRE_APPLY_TARGETS=template-stage template-log-level template-image-pull-policy
HELM_PRE_GENERATE_TARGETS = helm-values-update-image-version
HELM_POST_GENERATE_TARGETS = helm-values-replace-image-repo
CHECK_VAR_TARGETS=check-all-vars
IMAGE_IMPORT_TARGET=image-import
include build/make/k8s-component.mk

MOCKERY_IGNORED=vendor,build,docs,generated
include build/make/mocks.mk
include build/make/clean.mk
include makefiles/monitoring.mk
include makefiles/integration.mk

default: build

.PHONY: build
build: helm-delete image-import helm-apply ## Builds a new version of the k8s-ces-control and deploys it into the K8s-EcoSystem.

.PHONY: kill-pod
kill-pod:
	@echo "Restarting k8s-ces-control!"
	@kubectl -n ${NAMESPACE} delete pods -l 'app.kubernetes.io/name=k8s-ces-control'

.PHONY: helm-values-update-image-version
helm-values-update-image-version: $(BINARY_YQ)
	@echo "Updating the image version in source values.yaml to ${VERSION}..."
	@$(BINARY_YQ) -i e ".manager.image.tag = \"${VERSION}\"" ${K8S_COMPONENT_SOURCE_VALUES}

.PHONY: helm-values-replace-image-repo
helm-values-replace-image-repo: $(BINARY_YQ)
	@if [[ ${STAGE} == "development" ]]; then \
      		echo "Setting dev image repo in target values.yaml!" ;\
    		$(BINARY_YQ) -i e ".manager.image.registry=\"$(shell echo '${IMAGE_DEV}' | sed 's/\([^\/]*\)\/\(.*\)/\1/')\"" ${K8S_COMPONENT_TARGET_VALUES} ;\
    		$(BINARY_YQ) -i e ".manager.image.repository=\"$(shell echo '${IMAGE_DEV}' | sed 's/\([^\/]*\)\/\(.*\)/\2/')\"" ${K8S_COMPONENT_TARGET_VALUES} ;\
    	fi

.PHONY: template-stage
template-stage: $(BINARY_YQ)
	@if [[ ${STAGE} == "development" ]]; then \
  		echo "Setting STAGE env in deployment to ${STAGE}!" ;\
		$(BINARY_YQ) -i e ".manager.env.stage=\"${STAGE}\"" ${K8S_COMPONENT_TARGET_VALUES} ;\
	fi

.PHONY: template-log-level
template-log-level: ${BINARY_YQ}
	@if [[ "${STAGE}" == "development" ]]; then \
      echo "Setting LOG_LEVEL env in deployment to ${LOG_LEVEL}!" ; \
      $(BINARY_YQ) -i e ".manager.env.logLevel=\"${LOG_LEVEL}\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
    fi

.PHONY: template-image-pull-policy
template-image-pull-policy: $(BINARY_YQ)
	@if [[ "${STAGE}" == "development" ]]; then \
          echo "Setting pull policy to always!" ; \
          $(BINARY_YQ) -i e ".manager.imagePullPolicy=\"Always\"" "${K8S_COMPONENT_TARGET_VALUES}" ; \
    fi

.PHONY: check-env-var-stage
check-env-var-stage:
	@echo "Found stage [$(STAGE)]!"
	@$(call check_defined, STAGE, STAGE is not set. You need to export it before executing this command. Valid Values: [development, prodution])

.PHONY: check-env-var-log-level
check-env-var-log-level:
	@echo "Found log level [$(LOG_LEVEL)]!"
	@$(call check_defined, LOG_LEVEL, LOG_LEVEL is not set. You need to export it before executing this command. Valid Values: [DEBUG,INFO,WARN,ERROR])
