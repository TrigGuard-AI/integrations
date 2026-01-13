/**
 * TrigGuard JavaScript Demo
 * Minimal example demonstrating TG.execute() and guard()
 */

// Note: In a real project: import { execute, guard } from "@trigguard/sdk";

// Original payment function (before TrigGuard)
async function sendPayment(amount, recipient) {
  console.log(`💰 Sending ${amount} to ${recipient}`);
}

// Example 1: Using execute() directly
async function executeExample() {
  const receipt = await execute({
    action: "send_payment",
    payload: { amount: 100, recipient: "merchant@example.com" },
    context: {
      userId: "user_123",
      sessionId: "session_456",
      environment: "prod",
      requestId: "req_789"
    }
  });
  console.log(receipt.outcome === "executed" ? "✅ Payment executed" : "🔇 Payment not executed");
}

// Example 2: Using guard() wrapper (invisible integration)
async function guardExample() {
  const receipt = await guard(
    {
      action: "send_payment",
      payload: { amount: 100, recipient: "merchant@example.com" },
      context: { userId: "user_123", sessionId: "session_456", environment: "prod", requestId: "req_789" }
    },
    async () => { await sendPayment(100, "merchant@example.com"); }
  );
  console.log(`Outcome: ${receipt.outcome}`);
}
