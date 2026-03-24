# API mutation guard (example)

Illustrates calling a mutation endpoint protected by [`integrations/api-mutation-guard`](../integrations/api-mutation-guard/) middleware.

## Run the example server

```bash
cd integrations/api-mutation-guard
go run ./cmd/example-server
```

Server listens on **`:8080`**.

## Valid request shape

The issuer must have produced a receipt for **this exact body** and **surface** (`sid` in the payload). Headers carry the wire receipt and trust material.

```bash
curl -sS -X POST "http://localhost:8080/payments" \
  -H "Content-Type: application/json" \
  -H "TG-Execution-Receipt: <OER wire string>" \
  -H "TG-Surface: payment.execute" \
  -H "TG-Public-Key: <64 hex chars, Ed25519 public key>" \
  -d '{"amount":1000,"currency":"USD"}'
```

Use the same JSON bytes for `-d` that were hashed at issuance; otherwise **`ahsh`** will not match and verification fails.

## Failure behavior

- Missing or bad headers, unreadable body, signature / version / expiry / surface / hash mismatch → **HTTP 403** with body:

```text
OER authorization failed — mutation blocked
```

- Wrong HTTP method on example routes → **405 Method Not Allowed** (only after passing the guard for POST).

## Surfaces

Pick a surface id that matches your policy and receipt, e.g. `payment.execute`, `deploy.release`, or a custom id aligned with [EXECUTION_SURFACES.md](../docs/specs/EXECUTION_SURFACES.md).

See [docs/integrations/api_mutation_guard.md](../docs/integrations/api_mutation_guard.md).
