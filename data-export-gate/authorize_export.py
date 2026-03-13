#!/usr/bin/env python3
"""
TrigGuard data.export gate — call before running a data export.
Usage:
  export TRIGGUARD_URL TRIGGUARD_TOKEN
  python authorize_export.py --subject-digest "export_job_123" [--actor "data-service"]
Exits 0 on PERMIT, 1 on DENY or error. Writes receipt to receipt.json on success.
"""
import argparse
import json
import os
import sys
import urllib.request

def main():
    parser = argparse.ArgumentParser(description="TrigGuard data.export authorization")
    parser.add_argument("--subject-digest", required=True, help="Export job ID or dataset/scope hash")
    parser.add_argument("--actor", default="data-service", help="Actor ID (default: data-service)")
    parser.add_argument("--url", default=os.environ.get("TRIGGUARD_URL"), help="TrigGuard base URL")
    parser.add_argument("--token", default=os.environ.get("TRIGGUARD_TOKEN"), help="Bearer token")
    args = parser.parse_args()

    if not args.url or not args.token:
        print("Set TRIGGUARD_URL and TRIGGUARD_TOKEN", file=sys.stderr)
        sys.exit(1)

    body = json.dumps({
        "surface": "data.export",
        "actorId": args.actor,
        "subjectDigest": args.subject_digest,
    }).encode("utf-8")

    req = urllib.request.Request(
        f"{args.url.rstrip('/')}/execute",
        data=body,
        headers={
            "Authorization": f"Bearer {args.token}",
            "Content-Type": "application/json",
        },
        method="POST",
    )

    try:
        with urllib.request.urlopen(req) as resp:
            data = json.loads(resp.read().decode())
    except urllib.error.HTTPError as e:
        print(f"TrigGuard HTTP {e.code}: {e.read().decode()}", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"TrigGuard request failed: {e}", file=sys.stderr)
        sys.exit(1)

    decision = data.get("decision")
    if decision != "PERMIT":
        print(f"TrigGuard decision: {decision} (expected PERMIT)", file=sys.stderr)
        sys.exit(1)

    receipt = data.get("receipt", {})
    with open("receipt.json", "w") as f:
        json.dump(receipt, f, indent=2)
    print("TrigGuard PERMIT — receipt written to receipt.json")
    sys.exit(0)

if __name__ == "__main__":
    main()
