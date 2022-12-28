# Set these to the desired values
ARTIFACT_ID=k8s-ces-control
VERSION=0.0.1
GOTAG=1.19

# Setting SHELL to bash allows bash commands to be executed by recipes.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

## Image URL to use all building/pushing image targets
IMAGE_DEV=${K3CES_REGISTRY_URL_PREFIX}/${ARTIFACT_ID}:${VERSION}
IMAGE=cloudogu/${ARTIFACT_ID}:${VERSION}
LINT_VERSION=v1.45.2

MAKEFILES_VERSION=7.0.1
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

default: build

MONITORING_NAMESPACE=monitoring

.PHONY: expose-grafana
expose-grafana:
	@kubectl --namespace=${MONITORING_NAMESPACE} port-forward svc/grafana 8080:80 &

.PHONY: install-loki-stack
install-loki-stack: install-promtail install-minio install-loki install-grafana

.PHONY: delete-loki-stack
delete-loki-stack: delete-promtail delete-minio delete-loki delete-grafana

.PHONY: monitoring-namespace
monitoring-namespace:
	@kubectl create namespace monitoring || true

.PHONY: install-promtail
install-promtail: monitoring-namespace loki-example-credentials
	@kubectl apply -f ${WORKDIR}/loki-stack/promtail --namespace=${MONITORING_NAMESPACE}

.PHONY: loki-example-credentials
loki-example-credentials: monitoring-namespace
	@kubectl create secret generic loki-credentials --from-literal=username=loki --from-literal=password=loki  --namespace=${MONITORING_NAMESPACE} || true

.PHONY: delete-promtail
delete-promtail: monitoring-namespace
	@kubectl delete -f ${WORKDIR}/loki-stack/promtail --namespace=${MONITORING_NAMESPACE} || true

.PHONY: install-minio
install-minio: monitoring-namespace
	@kubectl apply -f ${WORKDIR}/loki-stack/minio --namespace=${MONITORING_NAMESPACE}

.PHONY: delete-minio
delete-minio: monitoring-namespace
	@kubectl delete -f ${WORKDIR}/loki-stack/minio --namespace=${MONITORING_NAMESPACE} || true

.PHONY: install-loki
install-loki: monitoring-namespace
	@kubectl apply -f ${WORKDIR}/loki-stack/loki/ --namespace=${MONITORING_NAMESPACE}
	@kubectl apply -f ${WORKDIR}/loki-stack/loki/write --namespace=${MONITORING_NAMESPACE}
	@kubectl apply -f ${WORKDIR}/loki-stack/loki/read --namespace=${MONITORING_NAMESPACE}
	@kubectl apply -f ${WORKDIR}/loki-stack/loki/gateway --namespace=${MONITORING_NAMESPACE}

.PHONY: delete-loki
delete-loki: monitoring-namespace
	@kubectl delete -f ${WORKDIR}/loki-stack/loki/ --namespace=${MONITORING_NAMESPACE} || true
	@kubectl delete -f ${WORKDIR}/loki-stack/loki/write --namespace=${MONITORING_NAMESPACE} || true
	@kubectl delete -f ${WORKDIR}/loki-stack/loki/read --namespace=${MONITORING_NAMESPACE} || true
	@kubectl delete -f ${WORKDIR}/loki-stack/loki/gateway --namespace=${MONITORING_NAMESPACE} || true

.PHONY: install-grafana
install-grafana: monitoring-namespace
	@kubectl apply -f ${WORKDIR}/loki-stack/grafana --namespace=${MONITORING_NAMESPACE}

.PHONY: delete-grafana
delete-grafana: monitoring-namespace
	@kubectl delete -f ${WORKDIR}/loki-stack/grafana --namespace=${MONITORING_NAMESPACE} || true

.PHONY: build
build: k8s-delete image-import k8s-apply ## Builds a new version of the k8s-ces-control and deploys it into the K8s-EcoSystem.

.PHONY: k8s-create-temporary-resource
k8s-create-temporary-resource: create-temporary-release-resource template-dev-only-image-pull-policy

.PHONY: create-temporary-release-resource
create-temporary-release-resource: $(K8S_RESOURCE_TEMP_FOLDER) check-env-var-stage check-env-var-log-level check-env-var-namespace
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

##@ GRpc-related Tasks

PROTOC_GEN_BIN=${UTILITY_BIN_PATH}/protoc-gen-go
PROTOC_GEN_VERSION=1.28.1
PROTOC_GEN_GRPC_BIN=${UTILITY_BIN_PATH}/protoc-gen-go-grpc
PROTOC_GEN_GRPC_VERSION=1.2
PROTOC_BUFFER_BIN=${UTILITY_BIN_PATH}/protoc
PROTOC_BUFFER_VERSION=21.12

#https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-linux-x86_64.zip

${PROTOC_GEN_BIN}: ${UTILITY_BIN_PATH}
	@echo Installing "protoc-gen-go"
	$(call go-get-tool,$(PROTOC_GEN_BIN),google.golang.org/protobuf/cmd/protoc-gen-go@v$(PROTOC_GEN_VERSION))

${PROTOC_GEN_GRPC_BIN}: ${UTILITY_BIN_PATH}
	@echo Installing "protoc-gen-go-grpc"
	$(call go-get-tool,$(PROTOC_GEN_GRPC_BIN),google.golang.org/grpc/cmd/protoc-gen-go-grpc@v$(PROTOC_GEN_GRPC_VERSION))

${PROTOC_BUFFER_BIN}: ${UTILITY_BIN_PATH}
	@echo Installing "protoc-buffer"
	@rm -rf $(UTILITY_BIN_PATH)/include
	@mkdir -p /tmp/protoc
	@wget -O /tmp/protoc/protoc-buffer.zip https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_BUFFER_VERSION)/protoc-$(PROTOC_BUFFER_VERSION)-linux-x86_64.zip
	@unzip -o /tmp/protoc/protoc-buffer.zip -d /tmp/protoc
	@mv -f /tmp/protoc/bin/protoc $(PROTOC_BUFFER_BIN)
	@mv -f /tmp/protoc/include $(UTILITY_BIN_PATH)
	@rm -rf /tmp/protoc

.PHONY install-grpc-tools:
install-grpc-tools: ${PROTOC_GEN_BIN} ${PROTOC_GEN_GRPC_BIN} ${PROTOC_BUFFER_BIN} ## Install all necessary gRPC tools

.PHONY generate-grpc:
generate-grpc: install-grpc-tools ## Generates the gRPC stubs for the protobuf files
	@echo Generating gRPC service code
	@rm -rf ./generated
	@mkdir generated
	@PATH="${PATH}:$(UTILITY_BIN_PATH)" $(PROTOC_BUFFER_BIN) \
--go_out=./${GENERATION_TARGET_DIR} \
--go_opt=module="github.com/cloudogu/k8s-ces-control/generated" \
--go-grpc_opt=module="github.com/cloudogu/k8s-ces-control/generated" \
--go-grpc_out=./${GENERATION_TARGET_DIR} \
-I ${GENERATION_SOURCE_DIR} \
./${GENERATION_SOURCE_DIR}/*.proto
	@git add ${GENERATION_TARGET_DIR}
	@echo "Make sure to update the generated mock files"