# Nginx execution gate example

## Start the stack

```bash
cd integrations/nginx-execution-gate
docker compose up --build
```

Nginx listens on `http://localhost:8080`.

## Send a gated request

This example expects the client to send the OER material as headers:

- `TG-Execution-Receipt`
- `TG-Surface`
- `TG-Public-Key`
- `TG-Action` — **base64** of the UTF-8 JSON body bytes (used as `action` for `ahsh` binding)

Example:

```bash
ACTION_JSON='{"amount":1000,"currency":"USD"}'
ACTION_B64=$(printf "%s" "$ACTION_JSON" | base64)

curl -i -X POST "http://localhost:8080/api/payments" \
  -H "Content-Type: application/json" \
  -H "TG-Action: ${ACTION_B64}" \
  -H "TG-Execution-Receipt: <receipt>" \
  -H "TG-Surface: payment.execute" \
  -H "TG-Public-Key: <64 hex chars>" \
  -d "$ACTION_JSON"
```

- If verification passes → Nginx proxies to `sample-app`.
- If verification fails → Nginx blocks the request (auth_request returns `403`).
