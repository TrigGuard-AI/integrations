#!/usr/bin/env node
/**
 * Protected target: simulates the irreversible action system.
 * ONLY accepts requests that come from the execution gateway (verified signature).
 * Path: /internal/commit — signature required; missing/invalid/replay → 403 UNAUTHORIZED_GATEWAY.
 * Path: /commit — legacy; x-trigguard-gateway: allowed → 200, else 403 DIRECT_EXECUTION_FORBIDDEN.
 */
const http = require('http');
const path = require('path');
const REPO_ROOT = path.resolve(__dirname, '../..');
const { verifyGatewayRequest, createNonceStore } = require(path.join(REPO_ROOT, 'src/core/verification/gatewaySignature.js'));

const PORT = Number(process.env.TARGET_PORT) || 3001;
const GATEWAY_HEADER = 'x-trigguard-gateway';
const GATEWAY_SECRET = process.env.GATEWAY_SECRET || '';
const nonceStore = createNonceStore({ ttlMs: 300000 });

function readBody(req) {
  return new Promise((resolve, reject) => {
    let buf = '';
    req.on('data', (ch) => { buf += ch; });
    req.on('end', () => resolve(buf));
    req.on('error', reject);
  });
}

const server = http.createServer(async (req, res) => {
  res.setHeader('Content-Type', 'application/json');
  const pathname = req.url && req.url.split('?')[0];

  // Legacy path: simple header check
  if (req.method === 'POST' && pathname === '/commit') {
    if (req.headers[GATEWAY_HEADER] !== 'allowed') {
      res.writeHead(403);
      res.end(JSON.stringify({ error: 'DIRECT_EXECUTION_FORBIDDEN' }));
      return;
    }
    res.writeHead(200);
    res.end(JSON.stringify({ ok: true, message: 'Protected action executed' }));
    return;
  }

  // Authority path: gateway signature required
  if (req.method === 'POST' && pathname === '/internal/commit') {
    const body = await readBody(req);
    const sig = req.headers['x-trigguard-gateway-signature'];
    const ts = req.headers['x-trigguard-gateway-timestamp'];
    const nonce = req.headers['x-trigguard-gateway-nonce'];
    if (!sig || !ts || !nonce || !GATEWAY_SECRET) {
      res.writeHead(403);
      res.end(JSON.stringify({ error: 'UNAUTHORIZED_GATEWAY' }));
      return;
    }
    const timestamp = parseInt(ts, 10);
    if (Number.isNaN(timestamp)) {
      res.writeHead(403);
      res.end(JSON.stringify({ error: 'UNAUTHORIZED_GATEWAY' }));
      return;
    }
    const ok = verifyGatewayRequest(
      { method: 'POST', path: '/internal/commit', body, signature: sig, timestamp, nonce },
      GATEWAY_SECRET,
      { nonceStore, maxAgeMs: 15000 }
    );
    if (!ok) {
      res.writeHead(403);
      res.end(JSON.stringify({ error: 'UNAUTHORIZED_GATEWAY' }));
      return;
    }
    res.writeHead(200);
    res.end(JSON.stringify({ ok: true, message: 'Protected action executed' }));
    return;
  }

  res.writeHead(404);
  res.end(JSON.stringify({ error: 'Not Found' }));
});

server.listen(PORT, () => {
  console.error(`[protected_target] listening on ${PORT} (/internal/commit = gateway signature required)`);
});
