# Reference verifier — authorization artifact before execution

**Purpose (State B):** A **minimal, runnable** illustration of the enforcement rule in [`docs/security/EXECUTION_ENFORCEMENT_MODEL.md`](../../docs/security/EXECUTION_ENFORCEMENT_MODEL.md):

> Protected execution surfaces **refuse** actions that do not carry a **valid** TrigGuard authorization artifact (receipt / token), with **binding** to the proposed action.

This is **not** production code. It proves the **mechanical** pattern: `verify(artifact, proposedAction) → allow | refuse`.

**Related:** [`examples/execution-gateway-demo/`](../execution-gateway-demo/) (HTTP gateway + token flow). This folder is intentionally **tiny** and dependency-light for reviewers.

## Run

```bash
cd examples/reference-verifier
npm install
npm run demo
```

Requires Node 18+.

## What the script does

1. Defines a **proposed action** (surface + payload).
2. Builds a **minimal authorization artifact** (decision, surface, `requestHash`, timestamp) aligned with the binding idea in the execution protocol.
3. Runs **`verifyExecution(artifact, proposed)`** — same checks a surface would apply before performing the side effect.
4. Prints **allowed** vs **refused** cases: missing artifact, wrong hash, expired TTL, decision not PERMIT.

## Integration

Real deployments should use the **canonical receipt shape** and **Ed25519 verification** from [`docs/protocol/TRIGGUARD_RECEIPT_SCHEMA.md`](../../docs/protocol/TRIGGUARD_RECEIPT_SCHEMA.md) and [`packages/trigguard-receipt-verifier`](../../packages/trigguard-receipt-verifier/) (or your language port). This example uses SHA-256 over a canonical payload string only to show **binding**, not to replace cryptographic verification.
