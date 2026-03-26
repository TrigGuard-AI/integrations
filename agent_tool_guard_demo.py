#!/usr/bin/env python3
"""
Minimal "agent → TrigGuard gate → PERMIT/DENY → (simulated) tool" demo.

Calls the canonical Swift decision service (POST /decide) — same authority path as
remote-eval-stub. No third-party deps (stdlib only).

Prerequisites
  1. Start the decision service (from repo root):

       ./scripts/start_canonical_core.sh

     (Exposes TG_CANONICAL_CORE_URL, typically http://127.0.0.1:9090/decide)

  2. Run this demo:

       python3 examples/agent_tool_guard_demo.py

  Override URL:

       TG_DECIDE_URL=http://127.0.0.1:9090/decide python3 examples/agent_tool_guard_demo.py

This is a credibility demo: it shows interception *before* a side-effect, not monitoring after.
"""
from __future__ import annotations

import json
import os
import sys
import urllib.error
import urllib.request

DEFAULT_DECIDE = os.environ.get("TG_CANONICAL_CORE_URL", "http://127.0.0.1:9090/decide")

# Align with remote-eval-stub/schemas.js (surface + signals + context shape for canonical invoke).
def spend_commit_frame(amount: float) -> dict[str, object]:
    return {
        "surface": "spendCommit",
        "signals": {
            "riskScore": min(1.0, amount / 1000.0),
            "dopamineDeficit": False,
            "rsdSpike": False,
        },
        "context": {
            "amount": amount,
            "origin": "automation",
        },
    }


def post_decide(url: str, body: dict[str, object]) -> tuple[int, dict | None]:
    data = json.dumps(body).encode("utf-8")
    req = urllib.request.Request(
        url,
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


def run_tool(name: str, amount: float) -> None:
    print(f"\n--- Agent requests tool: {name!r} (amount={amount}) ---")
    frame = spend_commit_frame(amount)
    try:
        status, resp = post_decide(DEFAULT_DECIDE, frame)
    except ConnectionError as e:
        print(f"Cannot reach decision service at {DEFAULT_DECIDE!r}: {e}")
        print("Start it from repo root: ./scripts/start_canonical_core.sh")
        sys.exit(2)
    if status != 200 or not resp:
        print(f"HTTP {status}: {resp}")
        sys.exit(1)
    decision = str(resp.get("decision", "")).upper()
    reason = resp.get("reasonCode", "")
    print(f"TrigGuard decision: {decision}  (reasonCode={reason!r})")
    if decision == "PERMIT":
        print("→ Tool execution ALLOWED (simulated: payment would proceed).")
    else:
        print("→ Tool execution BLOCKED (simulated: no side effect).")


def main() -> None:
    print("TrigGuard agent tool guard demo (canonical POST /decide)")
    print(f"Endpoint: {DEFAULT_DECIDE}")
    print("Scenario: two spendCommit attempts — policy may differ by amount.")
    run_tool("send_money", 500.0)
    run_tool("send_money", 50.0)
    print("\nDone.")


if __name__ == "__main__":
    main()
