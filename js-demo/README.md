# JavaScript Demo Example

**Minimal example demonstrating TrigGuard integration**

This example shows how to use `execute()` and `guard()` in a Node.js script.

## What It Demonstrates

- Using `execute()` directly
- Using `guard()` wrapper for invisible integration
- Handling receipt outcomes (executed vs silent)
- Silence is not an error

## Usage

```bash
# This is a reference example, not a runnable script
# Copy the pattern into your Node.js project
# npm install @trigguard/sdk
```

## Key Patterns

**Direct execution:**
```javascript
const receipt = await execute({
  action: "send_payment",
  payload: { amount: 100 },
  context: { ... }
});
```

**Guarded wrapper:**
```javascript
const receipt = await guard(
  { action: "send_payment", payload, context },
  async () => {
    await sendPayment(100, "merchant@example.com");
  }
);
```

The function only executes if TrigGuard returns "executed". If outcome is "silent", the function is not called (expected, not an error).

---

**Lines:** 58  
**Status:** Reference Example
