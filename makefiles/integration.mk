KUBECTL_BIN?=kubectl
GRPCURL_BIN?=${UTILITY_BIN_PATH}/grpcurl
GRPCURL_BIN_VERSION?=1.8.7

.PHONY: integration-test-bash
integration-test-bash: integration-test-bash-notice ${GRPCURL_BIN} ## Runs integration tests by bash.
	export GRPCURL_BIN=${GRPCURL_BIN} && export KUBECTL_BIN=${KUBECTL_BIN} && ./integration-test.sh

.PHONY: integration-test-bash-notice
integration-test-bash-notice:
	@echo "To run the integration tests be sure to have a local cluster set up and having the dogus postfix and ldap installed."
	@echo "The kubectl binary can be overwritten with KUBECTL_BIN."

${GRPCURL_BIN}: ${UTILITY_BIN_PATH}
	@echo "Installing grpcurl"
	$(call go-get-tool,$(GRPCURL_BIN),github.com/fullstorydev/grpcurl/cmd/grpcurl@v$(GRPCURL_BIN_VERSION))