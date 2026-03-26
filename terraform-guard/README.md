# Terraform execution guard (OER v1)

Offline verification of an **OER** for the **same** Terraform plan JSON bytes that were hashed at issuance (`ahsh`), then **`terraform apply`** only if verification succeeds.

Uses [`tools/oer-verifier-go`](../../tools/oer-verifier-go) — no network I/O in the verifier.

## Environment

| Variable | Required | Description |
|----------|----------|-------------|
| `TG_EXECUTION_RECEIPT` | yes | OER wire string |
| `TG_PUBLIC_KEY_HEX` | yes | Ed25519 public key (64 hex chars) |
| `TG_SURFACE` | yes | Must match payload `sid` (e.g. `infra.apply`) |
| `TG_ACTION_FILE` | no | Path to plan JSON; default **`terraform-plan.json`** |
| `TG_NOW_UNIX` | no | Clock for `exp` (default: now) |

Exit **0** → allow apply; **1** → block.

## Wrapper

[`terraform-apply-guard.sh`](terraform-apply-guard.sh) runs `go run .` then `terraform apply "$@"` if verification succeeds.

```bash
chmod +x terraform-apply-guard.sh
export TG_EXECUTION_RECEIPT='...'
export TG_PUBLIC_KEY_HEX='...'
export TG_SURFACE=infra.apply
export TG_ACTION_FILE=terraform-plan.json
./terraform-apply-guard.sh plan.tfplan
```

See [`examples/terraform-guard-example.md`](../../examples/terraform-guard-example.md).
