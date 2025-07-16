#!/bin/bash
set -eEuo pipefail
trap '__error_handing__ $? $LINENO' ERR

# This file is responsible to test the k8s-ces-control feature in integration with the whole cluster.
# To run this script a local cluster is needed.
# It is also required to install the following dogus: ldap, postfix

KUBECTL_BIN_PATH="${KUBECTL_BIN:-kubectl}"
GRPCURL_BIN_PATH="${GRPCURL_BIN:-grpcurl}"
GRPCURL_BIN_PATH_WITH_AUTH=
JQ_BIN_PATH="${JQ_BIN:-jq}"
PORT_FORWARD_PID=

function __error_handling__() {
  killPortForward
  local last_status_code=$1;
  local error_line_number=$2;
  echo 1>&2 "Error - exited with status $last_status_code at line $error_line_number";
  perl -slne 'if($.+5 >= $ln && $.-4 <= $ln){ $_="$. $_"; s/$ln/">" x length($ln)/eg; s/^\D+.*?$/\e[1;31m$&\e[0m/g;  print}' -- -ln=$error_line_number $0
}

startPortForward() {
  echo "Starting Port-Forward on ${PORT_FORWARD_PID}..."
  GRPCURL_PORT="$(python3 -c 'import socket; s=socket.socket(); s.bind(("", 0)); print(s.getsockname()[1]); s.close()')"
  "${KUBECTL_BIN_PATH}" port-forward service/k8s-ces-control "${GRPCURL_PORT}":50051 >/dev/null 2>&1 &
  PORT_FORWARD_PID=$!
  sleep 2s
  echo "Started Port-Forward on ${GRPCURL_PORT}"
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
  echo "Test-Suite: Dogu-Health: Test GetAll, GetByNames, GetByName"
  testDoguHealth_GetAll
  testDoguHealth_GetByNames
  testDoguHealth_GetByName
  echo "- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -"
  echo "Test-Suite: Support-Archive: Test Create"
  testSupportArchive_Create
  echo "- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -"
}

testDoguAdministration_GetAllDogus() {
  local doguListJson
  doguListJson=$(${GRPCURL_BIN_PATH} -plaintext localhost:"${GRPCURL_PORT}" doguAdministration.DoguAdministration.GetDoguList)
  local getDoguList
  getDoguList=$(echo ${doguListJson} | ${JQ_BIN_PATH} '.dogus | map(select(.name)) | .[].name')

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
  ${GRPCURL_BIN_PATH} -plaintext -d '{"doguName": "postfix"}' localhost:"${GRPCURL_PORT}" doguAdministration.DoguAdministration.StopDogu >/dev/null 2>&1

  # Wait for dogu to be terminated
  sleep 5s

  local replicas=""
  replicas="$(${KUBECTL_BIN_PATH} get deployment/postfix -o json | ${JQ_BIN_PATH} '.spec.replicas')"
  if [[ "${replicas}" == 0 ]]; then
    echo "Test: Postfix stopped? Success!"
    addSuccessTestCase "Dogu-Administration-StopDogu-Postfix" "k8s-ces-control successfully stopped the Postfix dogu."
  else
    echo "Test: Postfix stopped? Failed!"
    addFailingTestCase "Dogu-Administration-StopDogu-Postfix" "Expected the replicas of postfix to be 0 but got: ${replicas}"
  fi

  ${GRPCURL_BIN_PATH} -plaintext -d '{"doguName": "postfix"}' localhost:"${GRPCURL_PORT}" doguAdministration.DoguAdministration.StartDogu >/dev/null 2>&1

  # Wait for dogu to be started
  sleep 15s

  local healthStatus=""
  healthStatus="$(${KUBECTL_BIN_PATH} get dogu/postfix -o json | ${JQ_BIN_PATH} -r '.status.health')"
  if [[ "${healthStatus}" == "available" ]]; then
    echo "Test: Postfix started? Success!"
    addSuccessTestCase "Dogu-Administration-StartDogu-Postfix" "k8s-ces-control successfully started the Postfix dogu."
  else
    echo "Test: Postfix started? Failed!"
    addFailingTestCase "Dogu-Administration-StartDogu-Postfix" "Expected the healthStatus of postfix to be 'available' but got: ${healthStatus}"
  fi
}
testDoguAdministration_RestartDogus() {
  local postfixBeforePodName=""
  postfixBeforePodName="$(${KUBECTL_BIN_PATH} get pods -o name | grep postfix)"

  ${GRPCURL_BIN_PATH} -plaintext -d '{"doguName": "postfix"}' localhost:"${GRPCURL_PORT}" doguAdministration.DoguAdministration.RestartDogu >/dev/null 2>&1

  # Wait for dogu to be restarted
  sleep 15s

  local postfixAfterPodName=""
  postfixAfterPodName="$(${KUBECTL_BIN_PATH} get pods -o name | grep postfix)"

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

testDoguHealth_GetAll() {
  local allDogusHealthJson
  allDogusHealthJson=$(${GRPCURL_BIN_PATH} -plaintext localhost:"${GRPCURL_PORT}" health.DoguHealth.GetAll)

  if [[ $(echo ${allDogusHealthJson} | ${JQ_BIN_PATH} -r '.results.ldap.fullName') == 'ldap' && $(echo ${allDogusHealthJson} | ${JQ_BIN_PATH} -r '.results.ldap.healthy') == 'true' ]]; then
    echo "Test: [Dogu-Health-GetAll] Check if Ldap is healthy: Success!"
    addSuccessTestCase "Dogu-Health-GetAll-Ldap" "List of returned Dogus contained a healthy 'ldap' dogu."
  else
    echo "Test: [Dogu-Health-GetAll] Check if Ldap is healthy: Failed!"
    addFailingTestCase "Dogu-Health-GetAll-Ldap" "Expected to get Dogu 'ldap' is healthy but got only: ${allDogusHealthJson}"
  fi

  if [[ $(echo ${allDogusHealthJson} | ${JQ_BIN_PATH} -r '.results."nginx-static".fullName') == 'nginx-static' && $(echo ${allDogusHealthJson} | ${JQ_BIN_PATH} -r '.results."nginx-static".healthy') == 'true' ]]; then
      echo "Test: [Dogu-Health-GetAll] Check if NginxStatic is healthy: Success!"
      addSuccessTestCase "Dogu-Health-GetAll-NginxStatic" "List of returned Dogus contained a healthy 'nginx-static' dogu."
    else
      echo "Test: [Dogu-Health-GetAll] Check if NginxStatic is healthy: Failed!"
      addFailingTestCase "Dogu-Health-GetAll-NginxStatic" "Expected to get Dogu 'nginx-static' is healthy but got only: ${allDogusHealthJson}"
    fi

  if [[ $(echo ${allDogusHealthJson} | ${JQ_BIN_PATH} -r '.results."nginx-ingress".fullName') == 'nginx-ingress' && $(echo ${allDogusHealthJson} | ${JQ_BIN_PATH} -r '.results."nginx-ingress".healthy') == 'true' ]]; then
    echo "Test: [Dogu-Health-GetAll] Check if NginxIngress is healthy: Success!"
    addSuccessTestCase "Dogu-Health-GetAll-NginxIngress" "List of returned Dogus contained a healthy 'nginx-ingress' dogu."
  else
    echo "Test: [Dogu-Health-GetAll] Check if NginxIngress is healthy: Failed!"
    addFailingTestCase "Dogu-Health-GetAll-NginxIngress" "Expected to get Dogu 'nginx-ingress' is healthy but got only: ${allDogusHealthJson}"
  fi
}

testDoguHealth_GetByNames() {
  local dogusHealthJson
  dogusHealthJson=$(${GRPCURL_BIN_PATH} -plaintext -d '{"dogus": ["nginx-static", "ldap"]}' localhost:"${GRPCURL_PORT}" health.DoguHealth.GetByNames)

  if [[ $(echo ${dogusHealthJson} | ${JQ_BIN_PATH} -r '.results.ldap.fullName') == 'ldap' && $(echo ${dogusHealthJson} | ${JQ_BIN_PATH} -r '.results.ldap.healthy') == 'true' ]]; then
    echo "Test: [Dogu-Health-GetByNames] Check if Ldap is healthy: Success!"
    addSuccessTestCase "Dogu-Health-GetByNames-Ldap" "List of returned Dogus contained a healthy 'ldap' dogu."
  else
    echo "Test: [Dogu-Health-GetByNames] Check if Ldap is healthy: Failed!"
    addFailingTestCase "Dogu-Health-GetByNames-Ldap" "Expected to get Dogu 'ldap' is healthy but got only: ${dogusHealthJson}"
  fi

  if [[ $(echo ${dogusHealthJson} | ${JQ_BIN_PATH} -r '.results."nginx-static".fullName') == 'nginx-static' && $(echo ${dogusHealthJson} | ${JQ_BIN_PATH} -r '.results."nginx-static".healthy') == 'true' ]]; then
    echo "Test: [Dogu-Health-GetByNames] Check if NginxStatic is healthy: Success!"
    addSuccessTestCase "Dogu-Health-GetByNames-NginxStatic" "List of returned Dogus contained a healthy 'nginx-static' dogu."
  else
    echo "Test: [Dogu-Health-GetByNames] Check if NginxStatic is healthy: Failed!"
    addFailingTestCase "Dogu-Health-GetByNames-NginxStatic" "Expected to get Dogu 'nginx-static' is healthy but got only: ${dogusHealthJson}"
  fi

  if [[ $(echo ${dogusHealthJson} | ${JQ_BIN_PATH} -r '.results | length') == '2' ]]; then
    echo "Test: [Dogu-Health-GetByNames] Check result-length: Success!"
    addSuccessTestCase "Dogu-Health-GetByNames-ResultLength" "List of returned Dogus has 2 items."
  else
    echo "Test: [Dogu-Health-GetByNames] Check result-length: Failed!"
    addFailingTestCase "Dogu-Health-GetByNames-ResultLength" "Expected to get 2 items in result but got only: ${dogusHealthJson}"
  fi
}

testDoguHealth_GetByName() {
  local doguHealthJson
  doguHealthJson=$(${GRPCURL_BIN_PATH} -plaintext -d '{"dogu_name": "nginx-static"}' localhost:"${GRPCURL_PORT}" health.DoguHealth.GetByName)

  if [[ $(echo ${doguHealthJson} | ${JQ_BIN_PATH} -r '.fullName') == 'nginx-static' && $(echo ${doguHealthJson} | ${JQ_BIN_PATH} -r '.healthy') == 'true' ]]; then
    echo "Test: [Dogu-Health-GetByName] Check if NginxStatic is healthy: Success!"
    addSuccessTestCase "Dogu-Health-GetByName-NginxStatic" "List of returned Dogus contained a healthy 'nginx-static' dogu."
  else
    echo "Test: [Dogu-Health-GetByName] Check if NginxStatic is healthy: Failed!"
    addFailingTestCase "Dogu-Health-GetByName-NginxStatic" "Expected to get Dogu 'nginx-static' is healthy but got only: ${doguHealthJson}"
  fi

  ${GRPCURL_BIN_PATH} -plaintext -d '{"doguName": "nginx-static"}' localhost:"${GRPCURL_PORT}" doguAdministration.DoguAdministration.StopDogu >/dev/null 2>&1
  # Wait for dogu to be terminated
  sleep 5s

  doguHealthJson=$(${GRPCURL_BIN_PATH} -plaintext -d '{"dogu_name": "nginx-static"}' localhost:"${GRPCURL_PORT}" health.DoguHealth.GetByName)

  if [[ $(echo ${doguHealthJson} | ${JQ_BIN_PATH} -r '.fullName') == 'nginx-static' && $(echo ${doguHealthJson} | ${JQ_BIN_PATH} -r 'has("healthy")') == 'false' ]]; then
    echo "Test: [Dogu-Health-GetByName] Check if NginxStatic is not healthy: Success!"
    addSuccessTestCase "Dogu-Health-GetByName-NginxStatic" "List of returned Dogus contained a not-healthy 'nginx-static' dogu."
  else
    echo "Test: [Dogu-Health-GetByName] Check if NginxStatic is not healthy: Failed!"
    addFailingTestCase "Dogu-Health-GetByName-NginxStatic" "Expected to get Dogu 'nginx-static' is not-healthy but got only: ${doguHealthJson}"
  fi


  ${GRPCURL_BIN_PATH} -plaintext -d '{"doguName": "nginx-static"}' localhost:"${GRPCURL_PORT}" doguAdministration.DoguAdministration.StartDogu >/dev/null 2>&1
  # Wait for dogu to be started
  sleep 20s

  doguHealthJson=$(${GRPCURL_BIN_PATH} -plaintext -d '{"dogu_name": "nginx-static"}' localhost:"${GRPCURL_PORT}" health.DoguHealth.GetByName)

  if [[ $(echo ${doguHealthJson} | ${JQ_BIN_PATH} -r '.fullName') == 'nginx-static' && $(echo ${doguHealthJson} | ${JQ_BIN_PATH} -r '.healthy') == 'true' ]]; then
    echo "Test: [Dogu-Health-GetByName] Check if NginxStatic is healthy: Success!"
    addSuccessTestCase "Dogu-Health-GetByName-NginxStatic" "List of returned Dogus contained a healthy 'nginx-static' dogu."
  else
    echo "Test: [Dogu-Health-GetByName] Check if NginxStatic is healthy: Failed!"
    addFailingTestCase "Dogu-Health-GetByName-NginxStatic" "Expected to get Dogu 'nginx-static' is healthy but got only: ${doguHealthJson}"
  fi
}

testSupportArchive_Create() {
  local existingSupportArchives=""
  existingSupportArchives="$(${KUBECTL_BIN_PATH} get supportarchive -o custom-columns=NAME:.metadata.name --no-headers 2>/dev/null || echo "")"

  local createSupportArchiveJson
  createSupportArchiveJson=$(${GRPCURL_BIN_PATH} -plaintext -d '{"common": {"excluded_contents": {}, "logging_config": {}}}' localhost:"${GRPCURL_PORT}" maintenance.SupportArchive.Create)
  downloadPathGrpc=$(echo ${createSupportArchiveJson} | ${JQ_BIN_PATH} -r '.data' | base64 --decode)

  local newSupportArchives=""
  newSupportArchives="$(${KUBECTL_BIN_PATH} get supportarchive -o custom-columns=NAME:.metadata.name --no-headers 2>/dev/null || echo "")"

  sleep 5s # wait for support-archive-CR

  for archive in $newSupportArchives; do
    if ! echo "$existingSupportArchives" | grep -q "$archive"; then
      newlyCapturedArchives+=("$archive")
    fi
  done

  # Fail if no new archives were created
  if [ ${#newlyCapturedArchives[@]} -eq 0 ]; then
    echo "Test: [Support-Archive-Create] Check if support archive is created: Failed!"
    addFailingTestCase "Support-Archive-Create" "No new support archive was created after API call."
    return
  fi

  downloadPathArchive="$(${KUBECTL_BIN_PATH} get supportarchive ${newlyCapturedArchives[0]} -o json | ${JQ_BIN_PATH} -r '.status.downloadPath')"
  if [[ "$downloadPathGrpc" == "$downloadPathArchive" ]]; then
    echo "Test: [Support-Archive-Create] Check if support archive is created: Success!"
    addSuccessTestCase "Support-Archive-Create" "List of returned SupportArchive contained a created one."
  else
    echo "Test: [Support-Archive-Create] Check if support archive is created: Failed!"
    addFailingTestCase "Support-Archive-Create" "There is no existing support archive: ${createSupportArchiveJson}"
  fi
}

echo "Using KUBECTL=${KUBECTL_BIN_PATH}"
echo "Using GRPCURL=${GRPCURL_BIN_PATH}"
echo "Using JQ=${JQ_BIN_PATH}"

createIntegrationTestFile
startPortForward
runTests
killPortForward
finishIntegrationTestFile
