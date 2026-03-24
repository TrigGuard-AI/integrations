# Terraform execution guard (example)

This shows how to gate **`terraform apply`** with an **OER v1** receipt that binds to the **exact** Terraform plan JSON used for `ahsh` at issuance.

## 1. Plan and export JSON

```bash
terraform plan -out plan.tfplan
terraform show -json plan.tfplan > terraform-plan.json
```

The file **`terraform-plan.json`** is the UTF-8 bytes that must match what the TrigGuard issuer hashed (canonical per [CANONICAL_HASHING.md](../docs/specs/CANONICAL_HASHING.md)).

## 2. Request a receipt

Call your TrigGuard issuer with:

- **surface** `infra.apply` (or your governed surface; must match OER `sid`)
- **action** = the same `terraform-plan.json` content (or its canonical hash pipeline)

Store the wire receipt in **`TG_EXECUTION_RECEIPT`** and the trusted Ed25519 public key in **`TG_PUBLIC_KEY_HEX`**.

## 3. Run the guard, then apply

From [`integrations/terraform-guard/`](../integrations/terraform-guard/):

```bash
export TG_EXECUTION_RECEIPT='<wire receipt>'
export TG_PUBLIC_KEY_HEX='<64 hex chars>'
export TG_SURFACE=infra.apply
export TG_ACTION_FILE=terraform-plan.json   # default; set if path differs

chmod +x terraform-apply-guard.sh
./terraform-apply-guard.sh plan.tfplan
```

If verification **fails**, the script prints **`OER verification failed — terraform apply blocked`** and exits **1** — **`terraform apply` does not run**.

If verification **succeeds**, the script runs **`terraform apply plan.tfplan`** (or whatever arguments you pass).

## Notes

- Do not change `terraform-plan.json` after the receipt is issued; `ahsh` will not match.
- Invalid, expired, or wrong-surface receipts must fail closed (exit 1).

See [docs/integrations/terraform.md](../docs/integrations/terraform.md).
