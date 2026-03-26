#!/usr/bin/env bash
# TrigGuard infra.apply gate — call before terraform apply (or Pulumi/K8s equivalent).
# Requires: TRIGGUARD_URL, TRIGGUARD_TOKEN, SUBJECT_DIGEST (e.g. plan hash).
set -e

: "${TRIGGUARD_URL:?Set TRIGGUARD_URL}"
: "${TRIGGUARD_TOKEN:?Set TRIGGUARD_TOKEN}"
: "${SUBJECT_DIGEST:?Set SUBJECT_DIGEST (e.g. hash of terraform plan)}"

RESP=$(curl -s -w "\n%{http_code}" -X POST "$TRIGGUARD_URL/execute" \
  -H "Authorization: Bearer $TRIGGUARD_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"surface\":\"infra.apply\",\"actorId\":\"terraform\",\"subjectDigest\":\"$SUBJECT_DIGEST\"}")

HTTP_CODE=$(echo "$RESP" | tail -n1)
BODY=$(echo "$RESP" | sed '$d')

if [ "$HTTP_CODE" != "200" ]; then
  echo "TrigGuard HTTP $HTTP_CODE: $BODY"
  exit 1
fi

DECISION=$(echo "$BODY" | jq -r '.decision')
if [ "$DECISION" != "PERMIT" ]; then
  echo "TrigGuard decision: $DECISION (expected PERMIT)"
  exit 1
fi

echo "TrigGuard PERMIT — infra change authorized"
echo "$BODY" | jq -c '.receipt' > receipt_infra.json
echo "Receipt saved to receipt_infra.json"
