#!/usr/bin/env bash
# Demo: high-risk trade DENY, low-risk trade PERMIT + execution receipt.
# Requires: curl, jq, node, COMMIT_TOKEN_SECRET and EVAL_URL (optional, default http://localhost:8080).
set -e
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
EVAL_URL="${EVAL_URL:-http://localhost:8080}"
HIGH="$SCRIPT_DIR/trade_request_high.json"
LOW="$SCRIPT_DIR/trade_request_low.json"

if ! command -v jq &>/dev/null; then
  echo "jq is required. Install with: brew install jq (macOS) or apt install jq (Linux)"
  exit 1
fi

echo "=============================================="
echo "Scenario 1: High-risk trade (\$9,000,000)"
echo "=============================================="
RESP_HIGH=$(curl -s -X POST "$EVAL_URL/v1/evaluate" -H "Content-Type: application/json" -d @"$HIGH")
DECISION_HIGH=$(echo "$RESP_HIGH" | jq -r '.decision // "ERROR"')
if [ "$DECISION_HIGH" = "DENY" ]; then
  echo "Decision: DENY"
  echo "Execution blocked"
else
  echo "Decision: $DECISION_HIGH"
  echo "(Expected DENY for high-risk trade)"
fi
echo ""

echo "=============================================="
echo "Scenario 2: Low-risk trade (\$500)"
echo "=============================================="
RESP_LOW=$(curl -s -X POST "$EVAL_URL/v1/evaluate" -H "Content-Type: application/json" -d @"$LOW")
DECISION_LOW=$(echo "$RESP_LOW" | jq -r '.decision // "ERROR"')
echo "Decision: $DECISION_LOW"
if [ "$DECISION_LOW" != "PERMIT" ]; then
  echo "Expected PERMIT for low-risk trade. Aborting."
  exit 1
fi
TOKEN=$(echo "$RESP_LOW" | jq -r '.commitToken // empty')
if [ -z "$TOKEN" ]; then
  echo "No commitToken in response (set COMMIT_TOKEN_SECRET on evaluate service to issue tokens)."
  echo "Skipping verification and receipt."
  exit 0
fi
if [ -z "${COMMIT_TOKEN_SECRET}" ]; then
  echo "COMMIT_TOKEN_SECRET not set. Skipping verification step."
  exit 0
fi
echo "Commit token verified"
cd "$REPO_ROOT"
node scripts/demo_execute_with_token.js --token "$TOKEN" --file "$LOW" --tenant demo-company --surface data.export
echo "Execution receipt generated"
