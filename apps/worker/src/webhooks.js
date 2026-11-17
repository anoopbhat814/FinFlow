'use strict';
const crypto = require('node:crypto');

function sign(secret, timestamp, body) {
  return crypto.createHmac('sha256', secret).update(`${timestamp}.${body}`).digest('hex');
}

async function deliverWebhook(url, secret, event) {
  const body = JSON.stringify(event);
  const ts = Math.floor(Date.now() / 1000);
  const res = await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'FinFlow-Timestamp': String(ts),
      'FinFlow-Signature': `v1=${sign(secret, ts, body)}`,
    },
    body,
    signal: AbortSignal.timeout(5000),
  });
  if (!res.ok) throw new Error(`webhook responded ${res.status}`);
}

module.exports = { deliverWebhook, sign };
