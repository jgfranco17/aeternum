#!/usr/bin/env bash

SERVICE_INFO_URL="${API_BASE_URL}/service-info"
echo "Fetching service information: ${SERVICE_INFO_URL}"
curl "${SERVICE_INFO_URL}" | jq .

TEST_EXEC_URL="${API_BASE_URL}/v0/tests/run"
echo "Running sample test execution: ${TEST_EXEC_URL}"
curl -X POST "${TEST_EXEC_URL}" \
    --header "Content-Type: application/json" \
    -d @sample/basic_request.json | jq .
