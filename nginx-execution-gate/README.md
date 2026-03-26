# Nginx reverse-proxy execution gate (example)

This example demonstrates **enforcement before the application server** using:

- **Nginx** as the edge reverse proxy
- a small **verifier-service** that calls `oer.Verify` (offline)
- a dummy **sample-app** backend

It is an integration surface example only. It does **not** modify OER protocol specs, the verifier implementation, conformance assets, or SafetyEngine.

## Architecture

Client
  │
  ▼
Nginx Execution Gate (auth_request)
  │
  ▼
Verifier Service (`oer.Verify`)
  │
  ├─ 200 → proxy to backend
  └─ 403 → block at gateway

## Important note about request bodies

Nginx `auth_request` subrequests do **not** include the original request body by default.
To keep this example minimal and still bind `ahsh` to action bytes, the client supplies the action JSON bytes separately via `TG-Action` (base64-encoded), which Nginx forwards to the verifier.

In a production gateway, you would typically structure the system so the action bytes (or their canonical hash) are available to the gate in a safe, explicit way.

## Run

```bash
cd integrations/nginx-execution-gate
docker compose up --build
```

See `docs/nginx_execution_gate.md` and `examples/nginx-execution-gate-example.md` within this directory.
