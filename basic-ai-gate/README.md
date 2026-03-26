# Basic AI gate — conceptual demo

This folder is **educational only**. It does not implement the full TrigGuard product or runtime. Use it to explain the **integration shape** to engineers and stakeholders.

## Story: payment attempt

1. **AI agent wants to send a payment**  
   The agent (or orchestration layer) forms an execution request: what action, on what surface, with what context (amount, tenant, risk signals, etc.).

2. **TrigGuard intercepts the request**  
   Your application does **not** call the bank/API directly as the first step. It calls TrigGuard (or a gateway in front of it) with the normalized execution envelope.

3. **Policy evaluation occurs**  
   TrigGuard runs **deterministic** policy evaluation: rules, obligations, and any configured gates relevant to that surface.

4. **PERMIT or DENY**  
   The engine returns a **decision**. If the action is not allowed, the integration must **not** perform the side effect (fail closed).

5. **Receipt generated**  
   A **receipt** (or chain of receipts) records the decision and enough identity to audit and, in a full deployment, verify cryptographically.

## What to build for a real integration

- Normalize agent intent into the **protocol** execution shape.
- Call the **evaluation path** your deployment exposes (library, sidecar, or cloud).
- Enforce **DENY** and **SILENCE** semantics in your executor; only execute on **PERMIT** when your architecture requires it.
- Persist and propagate **receipts** per your compliance needs.

For code and packages, start from the main repository’s protocol package and server docs — not from this README alone.
