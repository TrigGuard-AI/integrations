# Data export gate — data.export

Authorize data exports (customer data, dataset dumps, AI dataset exports) through TrigGuard. Call TrigGuard **before** running the export; on PERMIT, run the export and keep the receipt for compliance and audit.

## Flow

1. Export request (e.g. job ID, scope, dataset).
2. **TrigGuard authorization** — Call `POST /execute` with surface `data.export` and a scope or job ID as `subjectDigest`.
3. If PERMIT, run the export; store the receipt.

## Script

[authorize_export.py](authorize_export.py) calls TrigGuard and returns the decision and receipt. Use it in your export pipeline:

```bash
export TRIGGUARD_URL="https://your-trigguard.run.app"
export TRIGGUARD_TOKEN="your-token"
python authorize_export.py --subject-digest "export_job_123" --actor "data-service"
# If exit 0, proceed with export and optionally save receipt.json for audit.
```

## Request

**POST /execute**

```json
{
  "surface": "data.export",
  "actorId": "data-service",
  "subjectDigest": "dataset_hash_or_export_job_id"
}
```

See [SURFACE_USAGE_EXAMPLES.md](../../docs/examples/SURFACE_USAGE_EXAMPLES.md).

## Surface

- **Surface:** `data.export`
- **subjectDigest:** Dataset hash, export job ID, or scope identifier.
- [TRIGGUARD_SURFACES.md](../../docs/protocol/TRIGGUARD_SURFACES.md)
