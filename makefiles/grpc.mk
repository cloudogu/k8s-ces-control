##@ GRPC-related Tasks

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