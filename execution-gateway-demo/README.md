# Execution Gateway Demo

**No token, no action.** This demo shows TrigGuard enforcement at the execution boundary:

- **Direct execution is forbidden** — The protected target rejects requests that do not come through the gateway.
- **Execution requires a valid TrigGuard commit token** — The gateway verifies the token; invalid or missing token → EXECUTION_DENIED.
- **High-risk request is denied** — Evaluate returns DENY; no token; gateway denies execution.
- **Low-risk request is permitted with receipt** — Evaluate returns PERMIT + commitToken; gateway verifies, forwards to target, returns execution receipt.

## What the demo proves

1. **Direct bypass attempt** → Protected target returns `403 DIRECT_EXECUTION_FORBIDDEN` (no `x-trigguard-gateway: allowed` header).
2. **High-risk action** → TrigGuard evaluate returns DENY; executing through gateway without a valid token → `EXECUTION_DENIED`.
3. **Low-risk action** → Evaluate returns PERMIT + commitToken; gateway verifies token, forwards to target, returns success and execution receipt.

The application cannot skip TrigGuard. The target only accepts requests that come through the gateway with a valid authorization.

## Prerequisites

- Node.js 18+
- curl, jq
- TrigGuard evaluate service running (e.g. `remote-eval-stub`) with `COMMIT_TOKEN_SECRET` set

## Run instructions

**Terminal 1 — TrigGuard evaluate service (required for scenarios 2 and 3):**

```bash
COMMIT_TOKEN_SECRET=dev-secret node remote-eval-stub/server.js
```

**Terminal 2 — Run the demo (starts protected target and gateway, runs three scenarios, then exits):**

```bash
export EVAL_URL=http://localhost:8080
export COMMIT_TOKEN_SECRET=dev-secret
./examples/execution-gateway-demo/run_demo.sh
```

Or run target and gateway manually:

```bash
# Terminal 1: protected target (port 3001)
COMMIT_TOKEN_SECRET=dev-secret node examples/execution-gateway-demo/protected_target.js

# Terminal 2: gateway (port 3002)
COMMIT_TOKEN_SECRET=dev-secret node examples/execution-gateway-demo/gateway.js

# Terminal 3: run scenarios with curl (see run_demo.sh for exact requests)
```

## Expected output

```
[1] Direct bypass attempt
    Result: BLOCKED (403 DIRECT_EXECUTION_FORBIDDEN)

[2] High-risk action
    Decision: DENY
    Gateway result: EXECUTION_DENIED

[3] Low-risk action
    Decision: PERMIT
    Gateway result: EXECUTION PERMITTED
    Execution receipt generated

Demo complete.
```

## Files

- **protected_target.js** — Simulated irreversible action server; only accepts requests with `x-trigguard-gateway: allowed`.
- **gateway.js** — Verifies commit token via `requireExecutionAuthorization`, forwards to target, returns receipt.
- **run_demo.sh** — Starts target and gateway, runs three scenarios, cleans up.
- **high_risk_request.json** / **low_risk_request.json** — Payloads for evaluate and gateway.
