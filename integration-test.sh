#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

# This file is responsible to test the k8s-ces-control feature in integration with the whole cluster.
# To run this script a local cluster is needed.

KUBECTL_BIN=kubectl
GRPCURL_BIN=grpcurl
GRPCURL_PORT=38341
PORT_FORWARD_PID=

startPortForward() {
  echo "Starting Port-Forward on ${PORT_FORWARD_PID}..."
  GRPCURL_PORT="$(python3 -c 'import socket; s=socket.socket(); s.bind(("", 0)); print(s.getsockname()[1]); s.close()')"
  "${KUBECTL_BIN}" port-forward service/k8s-ces-control ${GRPCURL_PORT}:50051 &
  PORT_FORWARD_PID=$!
  sleep 2s
  echo "Started Port-Forward on ${PORT_FORWARD_PID}"
}
killPortForward() {
  echo "Stopping Port-Forward..."
  kill -kill "${PORT_FORWARD_PID}"
}

test() {
local getDoguList
getDoguList="$(${GRPCURL_BIN} -insecure localhost:"${GRPCURL_PORT}" doguAdministration.DoguAdministration.GetDoguList | jq '.dogus | map(select(.name)) | .[].name')"
printf "Result: \n %s" "${getDoguList}"
#    String installedDogus = grpcurl(grpcurlPort, "")
#    echo "Retrieve all Dogus from "
#
#    String[] expectedDogus = ["ldap", "postfix"]
#    if (!installedDogus.contains("\"ldap\"")){
#        sh "echo 'Expected ldap dogu to be contained in the dogu list returned by grpc call but does not -> exit' && exit 1"
#    }
}

startPortForward
test
killPortForward
