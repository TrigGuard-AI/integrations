# Terraform infra gate — infra.apply

Gate Terraform apply (or other infra changes) with TrigGuard. Request authorization **before** running `terraform apply`; on PERMIT, run apply and keep the receipt for audit.

## Flow

1. `terraform plan` (or equivalent).
2. **TrigGuard authorization** — Call `POST /execute` with surface `infra.apply` and a plan hash or change ID as `subjectDigest`.
3. If PERMIT, run `terraform apply`; store the receipt.

## Script

[authorize_infra.sh](authorize_infra.sh) calls TrigGuard and exits 0 only on PERMIT. Use it in CI or locally:

```bash
export TRIGGUARD_URL="https://your-trigguard.run.app"
export TRIGGUARD_TOKEN="$(gcloud auth print-identity-token)"  # or your token
export SUBJECT_DIGEST="$(terraform plan -no-color -out=tfplan && sha256sum tfplan | cut -d' ' -f1)"
./authorize_infra.sh
terraform apply tfplan
```

## Request

**POST /execute**

```json
{
  "surface": "infra.apply",
  "actorId": "terraform",
  "subjectDigest": "terraform_plan_hash_or_change_id"
}
```

Decision and receipt are in the response. See [SURFACE_USAGE_EXAMPLES.md](../../docs/examples/SURFACE_USAGE_EXAMPLES.md).

## Surface

- **Surface:** `infra.apply`
- **subjectDigest:** Plan hash or change-set ID so the receipt binds to that specific change.
- [TRIGGUARD_SURFACES.md](../../docs/protocol/TRIGGUARD_SURFACES.md)
