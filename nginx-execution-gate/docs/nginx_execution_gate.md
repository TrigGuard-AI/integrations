# Reverse proxy execution gate (Nginx)

This integration example shows how TrigGuard enforcement can run **at the edge**, *before* application servers receive mutation requests.

## Architecture

Client
  │
  ▼
Nginx Execution Gate
  │
  ▼
Verifier Service (`oer.Verify`)
  │
  ▼
Application Backend

## What is verified

The verifier service calls `oer.Verify(receipt, actionBytes, surface, pubKey, now)` and therefore enforces:

- Ed25519 signature
- protocol version (`v`)
- expiry (`exp`)
- surface binding (`sid`)
- action hash binding (`ahsh`)

## Nginx flow

- Requests to `/api/` must pass `auth_request /verify`.
- Nginx performs an internal subrequest to `/verify`.
- `/verify` proxies to the verifier-service `/verify` endpoint.
- `200` allows proxying to the backend; `403` blocks the request.

## Notes

- The example uses the `TG-Action` header (base64 of UTF-8 JSON bytes) because Nginx `auth_request` subrequests do not include the original body by default.
