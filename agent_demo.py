#!/usr/bin/env python3
"""
Agent → TrigGuard (remote-eval) → PERMIT / DENY / SILENCE → (simulated) side effect.

Calls **POST /v1/evaluate** on `remote-eval-stub` — the same HTTP path production uses
(transport + audit + canonical core), not a direct POST to Swift /decide.

Prerequisites (two terminals + canonical core)

  Terminal A — Swift decision service:

    ./scripts/start_canonical_core.sh

  Terminal B — Node remote-eval (from repo root):

    cd remote-eval-stub
    export NODE_ENV=development
    export TRIGGUARD_UNSAFE_LOCAL_AUTH_BYPASS=true
    export TG_CANONICAL_CORE_URL=http://127.0.0.1:9090/decide
    node server.js

  Terminal C — this demo:

    python3 examples/agent_demo.py

  Override URL:

    TRIGGUARD_EVAL_URL=http://127.0.0.1:8080/v1/evaluate python3 examples/agent_demo.py

Policy semantics follow `remote-eval-stub/README.md` (e.g. HIGH_RISK_SPEND for large spend + high risk).
"""
from __future__ import annotations

import json
import os
import sys
import urllib.error
import urllib.request

EVAL_URL = os.environ.get("TRIGGUARD_EVAL_URL", "http://127.0.0.1:8080/v1/evaluate")


def evaluate(body: dict) -> tuple[int, dict | None]:
    data = json.dumps(body).encode("utf-8")
    req = urllib.request.Request(
        EVAL_URL,
        data=data,
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    try:
        with urllib.request.urlopen(req, timeout=15) as resp:
            return resp.status, json.loads(resp.read().decode("utf-8"))
    except urllib.error.HTTPError as e:
        raw = e.read().decode("utf-8", errors="replace")
        try:
            return e.code, json.loads(raw)
        except json.JSONDecodeError:
            return e.code, {"_raw": raw}
    except urllib.error.URLError as e:
        raise ConnectionError(str(e.reason or e)) from e


def run_scenario(name: str, body: dict) -> None:
    print(f"\n--- {name} ---")
    try:
        status, resp = evaluate(body)
    except ConnectionError as e:
        print(f"Cannot reach {EVAL_URL!r}: {e}")
        print("Start remote-eval + canonical core (see docstring at top of this file).")
        sys.exit(2)
    if status != 200 or not resp:
        print(f"HTTP {status}: {resp}")
        if status == 404:
            print("Hint: nothing served POST /v1/evaluate here — start remote-eval-stub (see docstring).")
        sys.exit(1)
    decision = str(resp.get("decision", "")).upper()
    rc = resp.get("reasonCode", "")
    print(f"TrigGuard decision: {decision}")
    print(f"reasonCode={rc!r}")
    if decision == "PERMIT":
        print("→ Simulated action: EXECUTE (e.g. payment / tool call proceeds).")
    else:
        print("→ Simulated action: BLOCKED (no irreversible side effect).")


def main() -> None:
    print("TrigGuard agent demo (remote-eval POST /v1/evaluate)")
    print(f"Endpoint: {EVAL_URL}")

    # Blocked path: matches remote-eval README "HIGH_RISK_SPEND" example.
    deny_spend = {
        "tenantId": "demo-tenant",
        "surface": "spendCommit",
        "signals": {"riskScore": 0.8, "dopamineDeficit": False, "rsdSpike": False},
        "context": {"amount": 5000, "origin": "human"},
    }
    run_scenario(
        "Agent requests large spend (policy yields DENY in this stub)",
        deny_spend,
    )

    # Allowed-style path: explicit surface from README (PERMIT).
    permit_time = {
        "tenantId": "demo-tenant",
        "surface": "timeCommit",
        "signals": {"riskScore": 0.3, "dopamineDeficit": False, "rsdSpike": False},
        "context": {"amount": 0, "origin": "human"},
    }
    run_scenario(
        "Agent requests timeCommit (low-risk path — expect PERMIT per stub rules)",
        permit_time,
    )

    print("\nDone.")


if __name__ == "__main__":
    main()
