#!/usr/bin/env node
/**
 * Protected target: simulates the irreversible action system.
 * ONLY accepts requests that come through the gateway (header x-trigguard-gateway: allowed).
 * Direct calls → 403 DIRECT_EXECUTION_FORBIDDEN.
 */
const http = require('http');

const PORT = Number(process.env.TARGET_PORT) || 3001;
const GATEWAY_HEADER = 'x-trigguard-gateway';

const server = http.createServer((req, res) => {
  res.setHeader('Content-Type', 'application/json');
  if (req.method !== 'POST' || (req.url && req.url.split('?')[0] !== '/commit')) {
    res.writeHead(404);
    res.end(JSON.stringify({ error: 'Not Found' }));
    return;
  }
  const allowed = req.headers[GATEWAY_HEADER];
  if (allowed !== 'allowed') {
    res.writeHead(403);
    res.end(JSON.stringify({ error: 'DIRECT_EXECUTION_FORBIDDEN' }));
    return;
  }
  res.writeHead(200);
  res.end(JSON.stringify({ ok: true, message: 'Protected action executed' }));
});

server.listen(PORT, () => {
  console.error(`[protected_target] listening on ${PORT} (only accepts x-trigguard-gateway: allowed)`);
});
