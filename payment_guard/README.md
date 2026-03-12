# Payment Guard Example

**Flow:** client → TrigGuard → payment execution → receipt → verifyReceipt()

Proves: external system → TrigGuard → irreversible action, with verifiable receipt.

## Run

1. Start TrigGuard (from repo root):

   ```bash
   TRIGGUARD_SECRET=test-secret node scripts/run_gateway_port.js
   ```
   (Default port 9340; set `PORT=8080` to use 8080.)

2. In another terminal (from repo root):

   ```bash
   BASE_URL=http://localhost:9340 TRIGGUARD_SECRET=test-secret node examples/payment_guard/client.js
   ```
   Or if gateway is on 8080: `BASE_URL=http://localhost:8080 node examples/payment_guard/client.js`

## Expected output

```
payment permitted
receipt verified
```

## What it does

- Builds an execution envelope (payments surface, commit token, nonce).
- POSTs to TrigGuard `/execute`.
- On success: builds a signed receipt from the response, verifies it with the receipt protocol (ed25519), then prints the two lines above.
