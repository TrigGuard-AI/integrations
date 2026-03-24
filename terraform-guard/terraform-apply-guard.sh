#!/usr/bin/env bash
# Run Terraform apply only after offline OER verification succeeds.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

if ! go run .; then
  echo "OER verification failed — terraform apply blocked" >&2
  exit 1
fi

exec terraform apply "$@"
