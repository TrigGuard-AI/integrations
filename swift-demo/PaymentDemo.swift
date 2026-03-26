import Foundation
import TrigGuardSDK

// Example: Wrapping a payment function with TrigGuard

// Original payment function (before TrigGuard)
func sendPayment(amount: Double, recipient: String) async {
    print("💰 Sending \(amount) to \(recipient)")
}

// Guarded payment function (after TrigGuard)
func guardedSendPayment(amount: Double, recipient: String, context: TGContext) async throws {
    let paymentData = try JSONEncoder().encode(["amount": amount, "recipient": recipient])
    
    // Wrap with TrigGuard - closure only executes if outcome is "executed"
    let receipt = try await TG.guard(
        action: "send_payment",
        payload: paymentData,
        context: context
    ) {
        await sendPayment(amount: amount, recipient: recipient)
    }
    
    // Receipt always returned (silence is not an error)
    switch receipt.outcome {
    case .executed: print("✅ Payment executed")
    case .silent: print("🔇 Payment not executed (expected, not an error)")
    }
}

// Example usage
func main() async throws {
    let context = TGContext(
        userId: "user_123",
        sessionId: "session_456",
        environment: .production,
        requestId: "req_789"
    )
    try await guardedSendPayment(amount: 100.0, recipient: "merchant@example.com", context: context)
}
