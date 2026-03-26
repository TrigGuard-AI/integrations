/**
 * Reference flow: proposed action → authorization artifact → verify → allow | refuse.
 * Aligns with docs/security/EXECUTION_ENFORCEMENT_MODEL.md (binding + refusal semantics).
 *
 * Run: npm run demo (from this directory)
 */

import { createHash } from "node:crypto";

/** Minimal proposed work (what the AI/client wants to do). */
export interface ProposedAction {
  surface: string;
  /** Normalized payload; in production, use the same canonicalization as the authority. */
  payload: Record<string, unknown>;
}

/**
 * Minimal artifact standing in for a receipt / permit token.
 * Production: use TRIGGUARD_RECEIPT_SCHEMA fields + Ed25519 signature verification.
 */
export interface AuthorizationArtifact {
  decision: "PERMIT" | "DENY" | "SILENCE";
  surface: string;
  /** SHA-256 hex over canonical payload bytes — binds artifact to this exact action body. */
  requestHash: string;
  /** ISO 8601 time of decision (used for TTL in this demo). */
  timestamp: string;
}

const TTL_SECONDS = 300;
/** Maximum allowed clock skew (seconds) — reject artifacts timestamped further in the future. */
const MAX_CLOCK_SKEW_SECONDS = 30;

function canonicalPayloadJson(payload: Record<string, unknown>): string {
  const keys = Object.keys(payload).sort();
  const obj: Record<string, unknown> = {};
  for (const k of keys) obj[k] = payload[k];
  return JSON.stringify(obj);
}

export function hashProposedPayload(payload: Record<string, unknown>): string {
  return createHash("sha256").update(canonicalPayloadJson(payload), "utf8").digest("hex");
}

export type VerifyResult =
  | { allowed: true }
  | { allowed: false; reason: string };

/**
 * Surface rule: refuse unless artifact is present, PERMIT, bound, and fresh.
 * Maps to: verify(artifact, proposed) in EXECUTION_ENFORCEMENT_MODEL.md
 */
export function verifyExecution(
  artifact: AuthorizationArtifact | null | undefined,
  proposed: ProposedAction,
  nowMs: number = Date.now()
): VerifyResult {
  if (artifact == null) {
    return { allowed: false, reason: "no authorization artifact" };
  }
  if (artifact.decision !== "PERMIT") {
    return { allowed: false, reason: `decision is ${artifact.decision}, not PERMIT` };
  }
  if (artifact.surface !== proposed.surface) {
    return { allowed: false, reason: "surface mismatch" };
  }
  const expectedHash = hashProposedPayload(proposed.payload);
  if (artifact.requestHash !== expectedHash) {
    return { allowed: false, reason: "requestHash does not bind to this payload (replay/wrong action)" };
  }
  const t = Date.parse(artifact.timestamp);
  if (Number.isNaN(t)) {
    return { allowed: false, reason: "invalid timestamp" };
  }
  if (t - nowMs > MAX_CLOCK_SKEW_SECONDS * 1000) {
    return { allowed: false, reason: "artifact timestamp is too far in the future" };
  }
  if (nowMs - t > TTL_SECONDS * 1000) {
    return { allowed: false, reason: `artifact expired (TTL ${TTL_SECONDS}s)` };
  }
  return { allowed: true };
}

function logCase(name: string, result: VerifyResult): void {
  const status = result.allowed ? "ALLOW execution" : "REFUSE execution";
  const detail = result.allowed ? "" : ` (${result.reason})`;
  console.log(`  [${name}] ${status}${detail}`);
}

function main(): void {
  const proposed: ProposedAction = {
    surface: "payment.execute",
    payload: { amount: 12000, currency: "USD" },
  };

  const goodHash = hashProposedPayload(proposed.payload);
  const goodArtifact: AuthorizationArtifact = {
    decision: "PERMIT",
    surface: "payment.execute",
    requestHash: goodHash,
    timestamp: new Date().toISOString(),
  };

  console.log("TrigGuard reference verifier — proposed action → artifact → verify\n");

  logCase("valid PERMIT + binding + fresh", verifyExecution(goodArtifact, proposed));

  logCase("no artifact", verifyExecution(undefined, proposed));

  const wrongSurface: AuthorizationArtifact = { ...goodArtifact, surface: "email.send" };
  logCase("wrong surface", verifyExecution(wrongSurface, proposed));

  const wrongHash: AuthorizationArtifact = {
    ...goodArtifact,
    requestHash: "0".repeat(64),
  };
  logCase("wrong requestHash (scope/replay)", verifyExecution(wrongHash, proposed));

  const denyArtifact: AuthorizationArtifact = { ...goodArtifact, decision: "DENY" };
  logCase("decision DENY", verifyExecution(denyArtifact, proposed));

  const old = new Date(Date.now() - (TTL_SECONDS + 10) * 1000).toISOString();
  const expired: AuthorizationArtifact = { ...goodArtifact, timestamp: old };
  logCase("expired timestamp", verifyExecution(expired, proposed));

  const future = new Date(Date.now() + 600_000).toISOString();
  const futureArtifact: AuthorizationArtifact = { ...goodArtifact, timestamp: future };
  logCase("future-dated timestamp", verifyExecution(futureArtifact, proposed));

  console.log("\nDone. Real systems add: Ed25519 verify, policy snapshot id, nonce/intent_id, executionId uniqueness.");
}

main();
