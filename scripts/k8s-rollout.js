#!/usr/bin/env node
'use strict';
// Usage: node scripts/k8s-rollout.js <namespace> <image-tag>
const { execFileSync } = require('node:child_process');

const [ns, tag] = process.argv.slice(2);
if (!ns || !tag) {
  console.error('usage: k8s-rollout.js <namespace> <image-tag>');
  process.exit(2);
}

const DEPLOYMENTS = { ledger: 'ghcr.io/finflow/ledger', fraud: 'ghcr.io/finflow/fraud', web: 'ghcr.io/finflow/web' };
const kubectl = (...args) => execFileSync('kubectl', ['-n', ns, ...args], { stdio: 'inherit' });

for (const [name, image] of Object.entries(DEPLOYMENTS)) {
  console.log(`==> ${name}: ${image}:${tag}`);
  kubectl('set', 'image', `deployment/${name}`, `${name}=${image}:${tag}`);
  try {
    kubectl('rollout', 'status', `deployment/${name}`, '--timeout=180s');
  } catch {
    console.error(`rollout of ${name} failed, undoing`);
    kubectl('rollout', 'undo', `deployment/${name}`);
    process.exit(1);
  }
}
