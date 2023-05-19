# Set these to the desired values
ARTIFACT_ID=k8s-ces-control
VERSION=0.0.1
GOTAG=1.19.3

# Setting SHELL to bash allows bash commands to be executed by recipes.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

## Image URL to use all building/pushing image targets
IMAGE_DEV?=${K3CES_REGISTRY_URL_PREFIX}/${ARTIFACT_ID}:${VERSION}
IMAGE?=cloudogu/${ARTIFACT_ID}:${VERSION}
LINT_VERSION=v1.50.1

MAKEFILES_VERSION=7.5.0
.DEFAULT_GOAL:=default
GENERATION_TARGET_DIR=generated
GENERATION_SOURCE_DIR=grpc-protobuf
INTEGRATION_TEST_NAME_PATTERN=.*_inttest$$
# make sure to create a statically linked binary
GO_BUILD_FLAGS=-mod=vendor -a -tags netgo,osusergo $(LDFLAGS) -o $(BINARY)

K8S_RESOURCE_DIR=${WORKDIR}/k8s
K8S_CES_CONTROL_RESOURCE_YAML=${K8S_RESOURCE_DIR}/k8s-ces-control.yaml
# set cluster root empty as we can ignore it for the k8s-ces-control
K8S_CLUSTER_ROOT=""

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
include build/make/k8s.mk
include makefiles/grpc.mk
include makefiles/monitoring.mk
include makefiles/integration.mk

default: build

.PHONY: build
build: check-env-var-namespace k8s-delete image-import k8s-apply kill-pod ## Builds a new version of the k8s-ces-control and deploys it into the K8s-EcoSystem.

.PHONY: kill-pod
kill-pod:
	@echo "Restarting k8s-ces-control!"
	@kubectl -n ${NAMESPACE} delete pods -l 'app.kubernetes.io/name=k8s-ces-control'

.PHONY: k8s-create-temporary-resource
k8s-create-temporary-resource: create-temporary-release-resource template-dev-only-image-pull-policy

.PHONY: create-temporary-release-resource
create-temporary-release-resource: $(K8S_RESOURCE_TEMP_FOLDER) check-env-var-stage check-env-var-log-level
	@echo "---" > $(K8S_RESOURCE_TEMP_YAML)
	@cat $(K8S_CES_CONTROL_RESOURCE_YAML) >> $(K8S_RESOURCE_TEMP_YAML)
	@sed -i "s/'{{\.LOG\_LEVEL}}'/$(LOG_LEVEL)/" $(K8S_RESOURCE_TEMP_YAML)
	@sed -i "s/'{{\.STAGE}}'/$(STAGE)/" $(K8S_RESOURCE_TEMP_YAML)
	@$(BINARY_YQ) -i e "(select(.kind == \"Deployment\").spec.template.spec.containers[]|select(.image == \"*$(ARTIFACT_ID)*\").image)=\"$(IMAGE)\"" $(K8S_RESOURCE_TEMP_YAML);

.PHONY: template-dev-only-image-pull-policy
template-dev-only-image-pull-policy: $(BINARY_YQ)
	@if [[ "${STAGE}""X" == "development""X" ]]; \
		then echo "Setting pull policy to always for development stage!" && $(BINARY_YQ) -i e "(select(.kind == \"Deployment\").spec.template.spec.containers[]|select(.image == \"*$(ARTIFACT_ID)*\").imagePullPolicy)=\"Always\"" $(K8S_RESOURCE_TEMP_YAML); \
	fi

STAGE?=production
.PHONY: check-env-var-stage
check-env-var-stage:
	@echo "Found stage [$(STAGE)]!"
	@$(call check_defined, STAGE, STAGE is not set. You need to export it before executing this command. Valid Values: [development, prodution])

LOG_LEVEL?=INFO
.PHONY: check-env-var-log-level
check-env-var-log-level:
	@echo "Found log level [$(LOG_LEVEL)]!"
	@$(call check_defined, LOG_LEVEL, LOG_LEVEL is not set. You need to export it before executing this command. Valid Values: [DEBUG,INFO,WARN,ERROR])

NAMESPACE?=ecosystem
.PHONY: check-env-var-namespace
check-env-var-namespace:
	@echo "Found namespace [$(NAMESPACE)]!"
	@$(call check_defined, NAMESPACE, NAMESPACE is not set. You need to export it before executing this command.)

### Reimplementation to also clean build/deb
.PHONY: clean
clean: $(ADDITIONAL_CLEAN) ## Remove target and tmp directories
	rm -rf ${TARGET_DIR}
	rm -rf ${TMP_DIR}
	rm -rf ${UTILITY_BIN_PATH}
	rm -rf build/deb

.PHONY: dist-clean
dist-clean: clean ## Remove all generated directories
	rm -rf node_modules
	rm -rf public/vendor
	rm -rf vendor
	rm -rf npm-cache
	rm -rf bower