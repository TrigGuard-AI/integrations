#!/usr/bin/env node
/**
 * Protected target: simulates the irreversible action system.
 * Only exposes POST /internal/commit. Verifies gateway identity via x-trigguard-gateway-signature (HMAC).
 * Direct or spoofed calls → 403 UNAUTHORIZED_GATEWAY or DIRECT_EXECUTION_FORBIDDEN.
 */
const http = require('http');
const path = require('path');
const REPO_ROOT = path.resolve(__dirname, '../..');
const { verifyGatewayRequest, createNonceStore } = require(path.join(REPO_ROOT, 'shared/utils/gatewaySignature.js'));

const PORT = Number(process.env.TARGET_PORT) || 3001;
const GATEWAY_SECRET = process.env.GATEWAY_SECRET || '';
const nonceStore = createNonceStore({ ttlMs: 300000 });

const server = http.createServer((req, res) => {
  res.setHeader('Content-Type', 'application/json');
  const pathname = req.url && req.url.split('?')[0];
  if (req.method !== 'POST' || pathname !== '/internal/commit') {
    res.writeHead(404);
    res.end(JSON.stringify({ error: 'Not Found' }));
    return;
  }

  let body = '';
  req.on('data', (chunk) => { body += chunk; });
  req.on('end', () => {
    if (!GATEWAY_SECRET) {
      res.writeHead(503);
      res.end(JSON.stringify({ error: 'GATEWAY_SECRET not configured' }));
      return;
    }
    const signature = req.headers['x-trigguard-gateway-signature'];
    const ts = req.headers['x-trigguard-gateway-timestamp'];
    const nonce = req.headers['x-trigguard-gateway-nonce'];
    const timestamp = ts != null ? Number.parseInt(ts, 10) : Number.NaN;
    const ok = Boolean(signature && Number.isInteger(timestamp) && nonce &&
      verifyGatewayRequest(
        { method: req.method, path: pathname, body, signature, timestamp, nonce },
        GATEWAY_SECRET,
        { maxAgeMs: 5000, nonceStore }
      ));
    if (!ok) {
      res.writeHead(403);
      res.end(JSON.stringify({ error: 'UNAUTHORIZED_GATEWAY' }));
      return;
    }
    res.writeHead(200);
    res.end(JSON.stringify({ ok: true, message: 'Protected action executed' }));
  });
});

server.listen(PORT, () => {
  console.error(`[protected_target] listening on ${PORT} (POST /internal/commit only; gateway signature required)`);
});
