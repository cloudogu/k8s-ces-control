#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

# This file is responsible to test the k8s-ces-control feature in integration with the whole cluster.
# To run this script a local cluster is needed.
# It is also required to install the following dogus: ldap, postfix

#KUBECTL_BIN=kubectl
#GRPCURL_BIN=grpcurl
GRPCURL_PORT=38341
PORT_FORWARD_PID=

startPortForward() {
  echo "Starting Port-Forward on ${PORT_FORWARD_PID}..."
  GRPCURL_PORT="$(python3 -c 'import socket; s=socket.socket(); s.bind(("", 0)); print(s.getsockname()[1]); s.close()')"
  "${KUBECTL_BIN}" port-forward service/k8s-ces-control "${GRPCURL_PORT}":50051 > /dev/null 2>&1 &
  PORT_FORWARD_PID=$!
  sleep 2s
  echo "Started Port-Forward on ${PORT_FORWARD_PID}"
}

killPortForward() {
  echo "Stopping Port-Forward..."
  kill -kill "${PORT_FORWARD_PID}" || true
}

INTEGRATION_TEST_RESULT_FOLDER=target/bash-integration-test
INTEGRATION_TEST_RESULT_FILE="${INTEGRATION_TEST_RESULT_FOLDER}"/results.xml
# Creates the xml test file containing the results for the tests.
function createIntegrationTestFile() {
  rm -rf "${INTEGRATION_TEST_RESULT_FOLDER}"
  mkdir -p "${INTEGRATION_TEST_RESULT_FOLDER}"
  touch "${INTEGRATION_TEST_RESULT_FILE}"

  cat <<EOT >"${INTEGRATION_TEST_RESULT_FILE}"
<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="github.com/cloudogu/k8s-ces-control/integration">
EOT
}

# Adds a successful test case to the result xml.
#
# Parameters:
# 1 - Name of the test case.
# 2 - Message of the test case.
function addSuccessTestCase() {
  local name="$1"
  local message="$2"

  cat <<EOT >>"${INTEGRATION_TEST_RESULT_FILE}"
  <testcase name="${name}">
    <system-out>${message}</system-out>
  </testcase>
EOT
}

# Adds a failing test case to the result xml.
#
# Parameters:
# 1 - Name of the test case.
# 2 - Message of the test case.
function addFailingTestCase() {
  local name="$1"
  local message="$2"

  cat <<EOT >>"${INTEGRATION_TEST_RESULT_FILE}"
  <testcase name="${name}">
    <failure message="bash integration test failed">
${message}
    </failure>
  </testcase>
EOT
}

# Finishes the xml syntax for the xml test file.
function finishIntegrationTestFile() {
  cat <<EOT >>"${INTEGRATION_TEST_RESULT_FILE}"
</testsuite>
EOT
}

runTests() {
  echo "- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -"
  echo "Test-Suite: Dogu-Administration: Test GetAllDogus, Start, Stop, Restart"
  testDoguAdministration_GetAllDogus
  testDoguAdministration_StartStopDogus
  testDoguAdministration_RestartDogus
  echo "- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -"
}

testDoguAdministration_GetAllDogus() {
  local getDoguList
  getDoguList="$(${GRPCURL_BIN} -insecure localhost:"${GRPCURL_PORT}" doguAdministration.DoguAdministration.GetDoguList | jq '.dogus | map(select(.name)) | .[].name')"

  if [[ "${getDoguList}" == *"\"ldap\""* ]]; then
    echo "Test: Check if Ldap returned: Success!"
    addSuccessTestCase "Dogu-Administration-GetAll-Ldap" "List of returned Dogus contained the 'ldap' dogu."
  else
    echo "Test: Check if Ldap returned: Failed!"
    addFailingTestCase "Dogu-Administration-GetAll-Ldap" "Expected to get Dogu 'ldap' but got only:

${getDoguList}"
  fi

  if [[ "${getDoguList}" == *"\"postfix\""* ]]; then
    echo "Test: Check if Postfix returned: Success!"
    addSuccessTestCase "Dogu-Administration-GetAll-Postfix" "List of returned Dogus contained the 'postfix' dogu."
  else
    echo "Test: Check if Postfix returned: Failed!"
    addFailingTestCase "Dogu-Administration-GetAll-Postfix" "Expected to get Dogu 'postfix' but got only:

${getDoguList}"
  fi
}
testDoguAdministration_StartStopDogus() {
  ${GRPCURL_BIN} -insecure -d '{"doguName": "postfix"}' localhost:"${GRPCURL_PORT}" doguAdministration.DoguAdministration.StopDogu > /dev/null 2>&1

  local replicas=""
  replicas="$(${KUBECTL_BIN} get deployment/postfix -o json | jq '.spec.replicas')"
  if [[ "${replicas}" == 0 ]]; then
    echo "Test: Postfix stopped? Success!"
    addSuccessTestCase "Dogu-Administration-StopDogu-Postfix" "k8s-ces-control successfully stopped the Postfix dogu."
  else
    echo "Test: Postfix stopped? Failed!"
    addFailingTestCase "Dogu-Administration-StopDogu-Postfix" "Expected the replicas of postfix to be 0 but got: ${replicas}"
  fi

  local replicas=""
  replicas="$(${KUBECTL_BIN} get deployment/postfix -o json | jq '.spec.replicas')"
  if [[ "${replicas}" == 0 ]]; then
    echo "Test: Postfix started? Success!"
    addSuccessTestCase "Dogu-Administration-StartDogu-Postfix" "k8s-ces-control successfully started the Postfix dogu."
  else
    echo "Test: Postfix started? Failed!"
    addFailingTestCase "Dogu-Administration-StartDogu-Postfix" "Expected the replicas of postfix to be 1 but got: ${replicas}"
  fi
}
testDoguAdministration_RestartDogus() {
  local postfixBeforePodName=""
  postfixBeforePodName="$(${KUBECTL_BIN} get pods -o name | grep postfix)"

  ${GRPCURL_BIN} -insecure -d '{"doguName": "postfix"}' localhost:"${GRPCURL_PORT}" doguAdministration.DoguAdministration.RestartDogu > /dev/null 2>&1

  local postfixAfterPodName=""
  postfixAfterPodName="$(${KUBECTL_BIN} get pods -o name | grep postfix)"

  if [[ "${postfixAfterPodName}" != "${postfixBeforePodName}" ]]; then
    echo "Test: Postfix restarted? Success!"
    addSuccessTestCase "Dogu-Administration-RestartDogu-Postfix" "k8s-ces-control successfully restarted Postfix."
  else
    echo "Test: Postfix restarted? Failed!"
    addFailingTestCase "Dogu-Administration-RestartDogu-Postfix" "Expected that k8s-ces-control restarted pod, but pod name did not change after restart request. However, it should do so as the restart consist of killing the old pod.
Name before killing the pod: ${postfixBeforePodName}
Name after killing the pod: ${postfixAfterPodName}"
  fi
}

echo "Using KUBECTL_BIN=${KUBECTL_BIN}"
echo "Using GRPCURL_BIN=${GRPCURL_BIN}"

createIntegrationTestFile
startPortForward
runTests
killPortForward
finishIntegrationTestFile
