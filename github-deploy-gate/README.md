# GitHub deploy gate — deploy.release

TrigGuard authorization gate for production deployments. Use this example to add a deploy gate to any GitHub Actions workflow.

## Flow

1. Build and test (your steps).
2. **TrigGuard authorization** — Call `POST /execute` with surface `deploy.release`. If decision is not PERMIT, the job fails.
3. Deploy (only runs after PERMIT).

## Setup

- Set `TRIGGUARD_URL` (repository variable or secret) to your TrigGuard Cloud URL, e.g. `https://trigguard-cloud-386138887132.europe-west2.run.app`.
- Set `TRIGGUARD_CLOUD_TOKEN` (secret) to an identity token or API token that can call the TrigGuard endpoint (required when the service uses authenticated access).

## Usage

Copy [.github/workflows/deploy.yml](.github/workflows/deploy.yml) into your repo, or add the TrigGuard step to your existing deploy workflow.

## Receipt

On PERMIT, the response includes a signed receipt. Store it for audit (e.g. as a workflow artifact or in `.receipts/trigguard/`). Verify offline with:

```bash
node scripts/verify_receipt.js receipt.json
```

## Surface

- **Surface:** `deploy.release`
- **subjectDigest:** Use `github.sha` (commit being deployed) so the receipt binds to that release.
- See [TRIGGUARD_SURFACES.md](../../docs/protocol/TRIGGUARD_SURFACES.md) and [SURFACE_USAGE_EXAMPLES.md](../../docs/examples/SURFACE_USAGE_EXAMPLES.md).
