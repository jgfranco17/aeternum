#!/usr/bin/env bash

if [[ -z "$API_BASE_URL" ]]; then
    echo "Please set API_BASE_URL to run this script!"
    exit 1
else
    echo "Using base URL: ${API_BASE_URL}"
fi

SERVICE_INFO_URL="${API_BASE_URL}/service-info"
echo "Fetching service information: ${SERVICE_INFO_URL}"
curl "${SERVICE_INFO_URL}" | jq .

TEST_EXEC_URL="${API_BASE_URL}/v0/tests/run"
echo "Running sample test execution: ${TEST_EXEC_URL}"
curl -X POST "${TEST_EXEC_URL}" \
    --header "Content-Type: application/json" \
    -d @sample/basic_request.json | jq .
