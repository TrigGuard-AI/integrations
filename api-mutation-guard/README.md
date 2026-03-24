# API mutation guard (OER v1)

HTTP **`Middleware`** that blocks handlers unless **`oer.Verify`** succeeds for:

- **`TG-Execution-Receipt`** — OER wire string  
- **`TG-Surface`** — must match payload `sid` (e.g. `payment.execute`)  
- **`TG-Public-Key`** — Ed25519 public key, **64 hex characters** (32 bytes)  

The **raw request body** is the action JSON whose hash must match receipt **`ahsh`** (same bytes the issuer hashed).

On failure: **HTTP 403** with body **`OER authorization failed — mutation blocked`**.

Uses [`tools/oer-verifier-go`](../../tools/oer-verifier-go) — no network I/O in the verifier.

## Example server

```bash
cd integrations/api-mutation-guard
go run ./cmd/example-server
```

See [`examples/api-guard-example.md`](../../examples/api-guard-example.md) and [`docs/integrations/api_mutation_guard.md`](../../docs/integrations/api_mutation_guard.md).
