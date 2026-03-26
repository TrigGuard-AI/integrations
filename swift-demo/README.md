# Swift Demo Example

**Minimal example demonstrating TrigGuard integration**

This example shows how to wrap an existing payment function with `TG.guard()` for nearly invisible integration.

## What It Demonstrates

- Wrapping existing functions with `TG.guard()`
- Handling receipt outcomes (executed vs silent)
- Silence is not an error
- Minimal integration (one wrapper call)

## Usage

```bash
# This is a reference example, not a runnable script
# Copy the pattern into your Swift project
```

## Key Pattern

```swift
let receipt = try await TG.guard(
    action: "send_payment",
    payload: paymentData,
    context: context
) {
    await sendPayment(amount: amount, recipient: recipient)
}
```

The closure only executes if TrigGuard returns "executed". If outcome is "silent", the closure is not called (expected, not an error).

---

**Lines:** 58  
**Status:** Reference Example
