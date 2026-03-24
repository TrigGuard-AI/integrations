# GitHub Action guard (OER v1)

Composite action that runs the **offline** OER verifier before irreversible steps (e.g. Terraform apply, releases).

- **No network I/O** in the verifier.
- Reads receipt from **`TG_EXECUTION_RECEIPT`** (or action input `receipt`).
- Binds **surface** (`sid`), **action JSON** file hash (`ahsh`), **exp**, **Ed25519 signature**, and **`v`** per [VERIFIER_SDK.md](../../docs/specs/VERIFIER_SDK.md).

## Usage

From a workflow in the same repository:

```yaml
- uses: ./integrations/github-action-guard
  with:
    receipt: ${{ secrets.TG_OER_RECEIPT }}
    surface: infra.apply
    action_file: plan.action.json
    public_key: ${{ vars.TG_OER_PUBLIC_KEY_HEX }}
```

`action_file` is resolved relative to the repository root (`GITHUB_WORKSPACE`) when set.

## Local / CI

```bash
export TG_EXECUTION_RECEIPT='<base64url(payload).base64url(sig)>'
export TG_SURFACE=infra.apply
export TG_ACTION_FILE=./plan.action.json
export TG_PUBLIC_KEY_HEX='<64 hex chars>'
go run .
```

Exit code **0** = verification success; **1** = failure (do not execute the irreversible step).
