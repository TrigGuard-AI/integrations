# TrigGuard Deploy Authorization Example

Example GitHub Actions step showing how TrigGuard authorizes a deployment and returns a signed receipt.

## Example step

```yaml
- name: Authorize deploy
  run: |
    curl -s -X POST https://trigguard.example.com/execute \
      -H "Content-Type: application/json" \
      -d '{
        "surface": "deploy.release",
        "nonce": "${{ github.run_id }}",
        "payload": {
          "service": "payments-api",
          "version": "${{ github.sha }}"
        }
      }'
```

(Replace `https://trigguard.example.com` with your TrigGuard gateway URL and include your envelope/auth as required.)

## Expected response

- **decision:** `PERMIT` | `BLOCK` | `SILENCE`
- **executionId:** unique execution identifier
- **receipt:** signed execution receipt (ed25519)

TrigGuard verifies the deployment request under the current policy snapshot and returns a signed execution receipt. The receipt can be verified later using the TrigGuard receipt verifier libraries (Node or Python) without trusting the server.

## Flow

```
CI pipeline
   │
   ▼
TrigGuard authorization (deploy.release)
   │
   ▼
Signed receipt
   │
   ▼
Independent verification (auditors, release logs)
```
