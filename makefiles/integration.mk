KUBECTL_BIN?=${UTILITY_BIN_PATH}/kubectl
KUBECTL_BIN_PATH?=${UTILITY_BIN_PATH}/kubectl
KUBECTL_BIN_VERSION=v1.26.0
GRPCURL_BIN?=${UTILITY_BIN_PATH}/grpcurl
GRPCURL_BIN_VERSION?=1.8.7
JQ_BIN?=${UTILITY_BIN_PATH}/jq
JQ_BIN_VERSION?=1.6

.PHONY: integration-test-bash
integration-test-bash: integration-test-bash-notice ${GRPCURL_BIN} ${JQ_BIN} ${KUBECTL_BIN}## Runs integration tests by bash.
	export GRPCURL_BIN=${GRPCURL_BIN} && export KUBECTL_BIN=${KUBECTL_BIN_PATH} && export JQ_BIN=${JQ_BIN} && ./integration-test.sh

.PHONY: integration-test-bash-notice
integration-test-bash-notice:
	@echo "To run the integration tests be sure to have a local cluster set up and having the dogus postfix and ldap installed."
	@echo "The kubectl binary can be overwritten with KUBECTL_BIN."

${GRPCURL_BIN}: ${UTILITY_BIN_PATH}
	@echo "Installing grpcurl v${GRPCURL_BIN_VERSION}"
	$(call go-get-tool,$(GRPCURL_BIN),github.com/fullstorydev/grpcurl/cmd/grpcurl@v$(GRPCURL_BIN_VERSION))

${JQ_BIN}: ${UTILITY_BIN_PATH}
	@echo "Installing jq v${JQ_BIN_VERSION}"
	@wget -O ${JQ_BIN} https://github.com/stedolan/jq/releases/download/jq-${JQ_BIN_VERSION}/jq-linux64
	@chmod +x ${JQ_BIN}

${KUBECTL_BIN}: ${UTILITY_BIN_PATH}
	@echo "Installing kubectl ${KUBECTL_BIN_VERSION}"
	@curl -o ${KUBECTL_BIN} -LO "https://dl.k8s.io/release/${KUBECTL_BIN_VERSION}/bin/linux/amd64/kubectl"
	@chmod +x ${KUBECTL_BIN}