# Trading Demo: Irreversible Execution

Two scenarios:

1. **High-risk trade ($9M)** → TrigGuard **DENY** → execution blocked.
2. **Low-risk trade ($500)** → TrigGuard **PERMIT** → commit token → verification → execution receipt.

## Prerequisites

- **curl**, **jq**, **Node.js** (18+)
- TrigGuard evaluate service running (e.g. `remote-eval-stub`) with `COMMIT_TOKEN_SECRET` set for token issuance
- `COMMIT_TOKEN_SECRET` in environment when running the script (for scenario 2 verification)

## Run locally

From repo root:

```bash
# Terminal 1: start evaluate service
COMMIT_TOKEN_SECRET=secret node remote-eval-stub/server.js

# Terminal 2: run demo
export EVAL_URL=http://localhost:8080
export COMMIT_TOKEN_SECRET=secret
./examples/demo-trading/run_demo.sh
```

Or from this directory:

```bash
EVAL_URL=http://localhost:8080 COMMIT_TOKEN_SECRET=secret ../run_demo.sh
```

(Use the path that reaches `run_demo.sh` from your cwd.)

## Expected output

```
==============================================
Scenario 1: High-risk trade ($9,000,000)
==============================================
Decision: DENY
Execution blocked

==============================================
Scenario 2: Low-risk trade ($500)
==============================================
Decision: PERMIT
Commit token verified
EXECUTION PERMITTED
...
EXECUTION RECEIPT
...
Execution receipt generated
```

## Files

- **trade_request_high.json** — High riskScore, dopamineDeficit, rsdSpike; large amount → DENY.
- **trade_request_low.json** — Low risk; small amount → PERMIT + token.
- **run_demo.sh** — Runs both scenarios (evaluate + optional verify/receipt).
