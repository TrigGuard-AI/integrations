#!/usr/bin/env node
/**
 * Payment guard example: client → TrigGuard → payment execution → receipt → verifyReceipt()
 * Run: BASE_URL=http://localhost:8080 node examples/payment_guard/client.js
 * Requires: TrigGuard server running; TRIGGUARD_SECRET (default: test-secret)
 */

const http = require('http');
const path = require('path');

const repoRoot = path.resolve(__dirname, '..', '..');
const baseUrl = process.env.BASE_URL || 'http://localhost:8080';
const secret = process.env.TRIGGUARD_SECRET || 'test-secret';

const { issueCommitToken, computeRequestHash } = require(path.join(repoRoot, 'src', 'core', 'capability', 'commitToken.js'));
const { normalizeRequestPayload } = require(path.join(repoRoot, 'src', 'core', 'verification', 'commitVerifier.js'));
const {
  createReceipt,
  signReceipt,
  verifyReceipt,
  generateKeyPair,
} = require(path.join(repoRoot, 'src', 'receipts', 'ExecutionReceiptProtocol.js'));

function post(url, body) {
  return new Promise((resolve, reject) => {
    const u = new URL(url);
    const data = JSON.stringify(body);
    const req = http.request({
      hostname: u.hostname,
      port: u.port || 80,
      path: u.pathname,
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'Content-Length': Buffer.byteLength(data) },
    }, (res) => {
      let buf = '';
      res.on('data', (c) => (buf += c));
      res.on('end', () => {
        try {
          resolve({ status: res.statusCode, body: buf ? JSON.parse(buf) : null });
        } catch {
          resolve({ status: res.statusCode, body: buf });
        }
      });
    });
    req.on('error', reject);
    req.write(data);
    req.end();
  });
}

async function main() {
  const payload = { surface: 'payments', signals: {}, context: { amount: 100, currency: 'USD' } };
  const requestHash = computeRequestHash(normalizeRequestPayload(payload));
  const nonce = `payment-guard-${Date.now()}-${Math.random().toString(36).slice(2)}`;
  const claims = {
    tenantId: 'T1',
    surface: 'payments',
    decision: 'PERMIT',
    requestHash,
    nonce,
    exp: Math.floor(Date.now() / 1000) + 300,
    policyVersion: '1',
    engineVersion: 'test',
  };
  const token = issueCommitToken(claims, secret);
  const envelope = {
    surface: 'payments',
    tenantId: 'T1',
    requestHash,
    nonce: claims.nonce,
    commitToken: token,
    issuedAt: Math.floor(Date.now() / 1000),
    expiresAt: Math.floor(Date.now() / 1000) + 300,
  };

  const res = await post(`${baseUrl}/execute`, { envelope, payload });

  if (res.status !== 200) {
    console.error('execute failed:', res.status, res.body);
    process.exit(1);
  }

  console.log('payment permitted');

  const executionId = res.body?.executionId ?? '';
  const decision = res.body?.decision ?? 'PERMIT';
  const policyVersion = res.body?.decisionHash ? '1' : '1';
  const timestamp = new Date().toISOString();

  const receipt = createReceipt({
    executionId,
    requestHash,
    surface: 'payments',
    decision,
    policyVersion,
    timestamp,
    issuer: 'payment-guard-example',
    publicKeyId: 'example-key-1',
  });

  const { publicKey, privateKey } = generateKeyPair();
  signReceipt(receipt, privateKey);

  const verification = verifyReceipt(receipt, publicKey);
  if (!verification.valid) {
    console.error('receipt verification failed:', verification.reason);
    process.exit(1);
  }

  console.log('receipt verified');
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});
