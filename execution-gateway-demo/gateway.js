#!/usr/bin/env node
/**
 * Execution gateway: verifies TrigGuard commit token before forwarding to protected target.
 * NO TOKEN -> NO EXECUTION. INVALID TOKEN -> NO EXECUTION. VALID TOKEN -> forward + receipt.
 */
const http = require('http');
const path = require('path');
const REPO_ROOT = path.resolve(__dirname, '../..');
const { requireExecutionAuthorization } = require(path.join(REPO_ROOT, 'shared/utils/commitVerifier.js'));
const { createExecutionReceipt } = require(path.join(REPO_ROOT, 'shared/utils/executionReceipt.js'));
const crypto = require('crypto');

const PORT = Number(process.env.GATEWAY_PORT) || 3002;
const TARGET_URL = process.env.TARGET_URL || 'http://localhost:3001';
const COMMIT_TOKEN_SECRET = process.env.COMMIT_TOKEN_SECRET || '';

function forwardToTarget(payload) {
  return new Promise((resolve, reject) => {
    const body = JSON.stringify(payload);
    const u = new URL(TARGET_URL + '/commit');
    const req = http.request({
      hostname: u.hostname,
      port: u.port || 80,
      path: u.pathname,
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Content-Length': Buffer.byteLength(body),
        'x-trigguard-gateway': 'allowed',
      },
    }, (res) => {
      let data = '';
      res.on('data', (chunk) => { data += chunk; });
      res.on('end', () => {
        try {
          const out = data ? JSON.parse(data) : {};
          resolve({ statusCode: res.statusCode, body: out });
        } catch (e) {
          resolve({ statusCode: res.statusCode, body: data });
        }
      });
    });
    req.on('error', reject);
    req.setTimeout(5000, () => { req.destroy(); reject(new Error('Target timeout')); });
    req.end(body);
  });
}

const server = http.createServer(async (req, res) => {
  res.setHeader('Content-Type', 'application/json');
  if (req.method !== 'POST' || (req.url && req.url.split('?')[0] !== '/execute')) {
    res.writeHead(404);
    res.end(JSON.stringify({ error: 'Not Found' }));
    return;
  }
  let body = '';
  for await (const chunk of req) body += chunk;
  let parsed;
  try {
    parsed = body ? JSON.parse(body) : {};
  } catch {
    res.writeHead(400);
    res.end(JSON.stringify({ ok: false, error: 'Invalid JSON' }));
    return;
  }
  const payload = parsed.payload || parsed;
  const commitToken = parsed.commitToken || req.headers['x-trigguard-commit-token'];

  if (!COMMIT_TOKEN_SECRET) {
    res.writeHead(500);
    res.end(JSON.stringify({ ok: false, error: 'COMMIT_TOKEN_SECRET not set' }));
    return;
  }

  if (!commitToken || typeof commitToken !== 'string' || commitToken.trim() === '') {
    res.writeHead(403);
    res.end(JSON.stringify({ ok: false, error: 'EXECUTION_DENIED', reason: 'missing or invalid token' }));
    return;
  }

  const requestPayload = {
    tenantId: payload.tenantId,
    surface: payload.surface || 'spendCommit',
    signals: payload.signals || {},
    context: payload.context || {},
  };

  try {
    const { claims } = await requireExecutionAuthorization({
      token: commitToken,
      requestPayload,
      expectedSurface: 'spendCommit',
      expectedTenantId: payload.tenantId,
      secret: COMMIT_TOKEN_SECRET,
    });

    const fwd = await forwardToTarget(payload);
    if (fwd.statusCode !== 200) {
      res.writeHead(502);
      res.end(JSON.stringify({ ok: false, error: 'Target rejected', targetStatus: fwd.statusCode }));
      return;
    }

    const executionId = crypto.randomUUID();
    const { receipt, receiptHash } = createExecutionReceipt({
      tokenClaims: claims,
      executionId,
      commitToken,
    });

    res.writeHead(200);
    res.end(JSON.stringify({
      ok: true,
      message: 'Execution permitted through gateway',
      receipt,
      receiptHash,
    }));
  } catch (e) {
    if (e.code === 'EXECUTION_DENIED') {
      res.writeHead(403);
      res.end(JSON.stringify({ ok: false, error: 'EXECUTION_DENIED', reason: e.reason || e.message }));
      return;
    }
    res.writeHead(500);
    res.end(JSON.stringify({ ok: false, error: 'Internal error', message: e.message }));
  }
});

server.listen(PORT, () => {
  console.error(`[gateway] listening on ${PORT} (TARGET_URL=${TARGET_URL})`);
});
