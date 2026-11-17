#!/usr/bin/env node
'use strict';
// Post-deploy smoke test: hits health endpoints of every service.
const targets = {
  ledger: process.env.LEDGER_URL || 'http://localhost:8080/healthz',
  fraud: process.env.FRAUD_SERVICE_URL ? `${process.env.FRAUD_SERVICE_URL}/healthz` : 'http://localhost:8000/healthz',
};

(async () => {
  let failed = 0;
  for (const [name, url] of Object.entries(targets)) {
    try {
      const res = await fetch(url, { signal: AbortSignal.timeout(3000) });
      console.log(`${res.ok ? 'PASS' : 'FAIL'} ${name} (${res.status})`);
      if (!res.ok) failed++;
    } catch (err) {
      console.log(`FAIL ${name}: ${err.message}`);
      failed++;
    }
  }
  process.exit(failed ? 1 : 0);
})();
