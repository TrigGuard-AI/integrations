#!/usr/bin/env node
/**
 * TrigGuard payment.execute gate — call before executing a critical payment.
 * Usage:
 *   TRIGGUARD_URL=... TRIGGUARD_TOKEN=... node authorize_payment.js --subject-digest "txn_123" [--actor "payments-api"]
 * Exits 0 on PERMIT, 1 on DENY or error. Writes receipt to receipt.json on success.
 */
const url = process.env.TRIGGUARD_URL;
const token = process.env.TRIGGUARD_TOKEN;

const args = process.argv.slice(2);
let subjectDigest = null;
let actor = "payments-api";
for (let i = 0; i < args.length; i++) {
  if (args[i] === "--subject-digest" && args[i + 1]) {
    subjectDigest = args[i + 1];
    i++;
  } else if (args[i] === "--actor" && args[i + 1]) {
    actor = args[i + 1];
    i++;
  }
}

if (!url || !token) {
  console.error("Set TRIGGUARD_URL and TRIGGUARD_TOKEN");
  process.exit(1);
}
if (!subjectDigest) {
  console.error("Pass --subject-digest <transaction_id_or_reference>");
  process.exit(1);
}

const body = JSON.stringify({
  surface: "payment.execute",
  actorId: actor,
  subjectDigest,
});

fetch(`${url.replace(/\/$/, "")}/execute`, {
  method: "POST",
  headers: {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  },
  body,
})
  .then((res) => {
    if (!res.ok) {
      return res.text().then((t) => {
        throw new Error(`HTTP ${res.status}: ${t}`);
      });
    }
    return res.json();
  })
  .then((data) => {
    if (data.decision !== "PERMIT") {
      console.error(`TrigGuard decision: ${data.decision} (expected PERMIT)`);
      process.exit(1);
    }
    const fs = require("fs");
    fs.writeFileSync("receipt.json", JSON.stringify(data.receipt, null, 2));
    console.log("TrigGuard PERMIT — receipt written to receipt.json");
    process.exit(0);
  })
  .catch((err) => {
    console.error("TrigGuard request failed:", err.message);
    process.exit(1);
  });
