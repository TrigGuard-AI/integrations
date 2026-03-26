#!/usr/bin/env bash
# Execution gateway demo: NO TOKEN -> NO EXECUTION. INVALID TOKEN -> NO EXECUTION. VALID TOKEN -> EXECUTION ALLOWED.
# Requires: curl, jq, node. Start TrigGuard evaluate service (e.g. remote-eval-stub) and set EVAL_URL, COMMIT_TOKEN_SECRET.
set -e
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
EVAL_URL="${EVAL_URL:-http://localhost:8080}"
TARGET_PORT="${TARGET_PORT:-3001}"
GATEWAY_PORT="${GATEWAY_PORT:-3002}"
HIGH="$SCRIPT_DIR/high_risk_request.json"
LOW="$SCRIPT_DIR/low_risk_request.json"

if ! command -v jq &>/dev/null; then
  echo "jq is required. Install with: brew install jq (macOS) or apt install jq (Linux)"
  exit 1
fi

# Start protected target and gateway in background
export TARGET_PORT GATEWAY_PORT TARGET_URL="http://localhost:$TARGET_PORT"
export COMMIT_TOKEN_SECRET="${COMMIT_TOKEN_SECRET:-dev-secret}"
export GATEWAY_SECRET="${GATEWAY_SECRET:-dev-gateway-secret}"

node "$SCRIPT_DIR/protected_target.js" &
PID_TARGET=$!
node "$SCRIPT_DIR/gateway.js" &
PID_GATEWAY=$!

cleanup() {
  kill $PID_TARGET $PID_GATEWAY 2>/dev/null || true
  wait $PID_TARGET $PID_GATEWAY 2>/dev/null || true
}
trap cleanup EXIT

sleep 2

echo ""
echo "[1] Direct execution attempt → BLOCKED"
echo "    Calling protected target directly (no gateway header)..."
RESP_DIRECT=$(curl -s -o /dev/null -w "%{http_code}" -X POST "http://localhost:$TARGET_PORT/internal/commit" -H "Content-Type: application/json" -d '{}')
if [ "$RESP_DIRECT" = "403" ]; then
  echo "    Result: BLOCKED (403 DIRECT_EXECUTION_FORBIDDEN)"
else
  echo "    Result: unexpected status $RESP_DIRECT (expected 403)"
fi
echo ""

echo "[2] High risk request → DENY → execution blocked"
echo "    Evaluate high-risk request..."
RESP_EVAL=$(curl -s -X POST "$EVAL_URL/v1/evaluate" -H "Content-Type: application/json" -d @"$HIGH")
DECISION=$(echo "$RESP_EVAL" | jq -r '.decision // "ERROR"')
echo "    Decision: $DECISION"
if [ "$DECISION" = "DENY" ]; then
  echo "    Attempting execute through gateway with no token..."
  GATEWAY_BODY=$(jq -n --argjson p "$(jq -c . "$HIGH")" '{payload: $p}')
  GATEWAY_RESP=$(curl -s -X POST "http://localhost:$GATEWAY_PORT/execute" -H "Content-Type: application/json" -d "$GATEWAY_BODY")
  ERR=$(echo "$GATEWAY_RESP" | jq -r '.error // ""')
  echo "    Gateway result: ${ERR:-EXECUTION_DENIED}"
else
  echo "    (Expected DENY for high-risk; gateway would still deny without valid token)"
fi
echo ""

echo "[3] Low risk request → PERMIT → execution allowed"
echo "    Evaluate low-risk request..."
RESP_LOW=$(curl -s -X POST "$EVAL_URL/v1/evaluate" -H "Content-Type: application/json" -d @"$LOW")
DECISION_LOW=$(echo "$RESP_LOW" | jq -r '.decision // "ERROR"')
echo "    Decision: $DECISION_LOW"
if [ "$DECISION_LOW" != "PERMIT" ]; then
  echo "    Expected PERMIT. Ensure evaluate service is running with COMMIT_TOKEN_SECRET."
  exit 1
fi
TOKEN=$(echo "$RESP_LOW" | jq -r '.commitToken // empty')
if [ -z "$TOKEN" ]; then
  echo "    No commitToken (set COMMIT_TOKEN_SECRET on evaluate service)."
  exit 1
fi
echo "    Execute through gateway with commit token..."
EXEC_BODY=$(jq -n --argjson p "$(jq -c . "$LOW")" --arg t "$TOKEN" '{payload: $p, commitToken: $t}')
EXEC_RESP=$(curl -s -X POST "http://localhost:$GATEWAY_PORT/execute" -H "Content-Type: application/json" -d "$EXEC_BODY")
OK=$(echo "$EXEC_RESP" | jq -r '.ok // false')
if [ "$OK" = "true" ]; then
  echo "    Gateway result: EXECUTION PERMITTED"
  echo "    Execution receipt generated"
else
  echo "    Gateway result: $(echo "$EXEC_RESP" | jq -r '.error // .message // "failed"')"
  exit 1
fi
echo ""
echo "Demo complete."
