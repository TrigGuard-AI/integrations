# Payment gate — payment.execute

Authorize critical payments or financial transactions through TrigGuard. Call TrigGuard **before** executing the payment; on PERMIT, execute and keep the receipt for auditors.

## Flow

1. Payment request (transaction ID, amount, reference).
2. **TrigGuard authorization** — Call `POST /execute` with surface `payment.execute` and transaction ID (or reference) as `subjectDigest`.
3. If PERMIT, execute the payment; store the receipt.

## Script

[authorize_payment.js](authorize_payment.js) calls TrigGuard and prints the decision; on PERMIT writes the receipt to `receipt.json`. Node 18+ (fetch built-in).

```bash
export TRIGGUARD_URL="https://your-trigguard.run.app"
export TRIGGUARD_TOKEN="your-token"
node authorize_payment.js --subject-digest "txn_abc123" [--actor "payments-api"]
# If exit 0, proceed with payment; receipt in receipt.json
```

## Request

**POST /execute**

```json
{
  "surface": "payment.execute",
  "actorId": "payments-api",
  "subjectDigest": "transaction_id_or_reference"
}
```

See [SURFACE_USAGE_EXAMPLES.md](../../docs/examples/SURFACE_USAGE_EXAMPLES.md).

## Surface

- **Surface:** `payment.execute`
- **subjectDigest:** Transaction ID or payment reference so the receipt binds to that transaction.
- [TRIGGUARD_SURFACES.md](../../docs/protocol/TRIGGUARD_SURFACES.md)
