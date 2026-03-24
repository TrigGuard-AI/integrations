#!/usr/bin/env bash
# Run Terraform apply only after offline OER verification succeeds.
set -euo pipefail

ORIG_DIR="$(pwd)"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

if ! TG_CALLER_DIR="$ORIG_DIR" go run .; then
  echo "OER verification failed — terraform apply blocked" >&2
  exit 1
fi

cd "$ORIG_DIR"
exec terraform apply "$@"
