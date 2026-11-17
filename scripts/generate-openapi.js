#!/usr/bin/env node
'use strict';
// Emits docs/openapi.json from a compact route table.
const fs = require('node:fs');
const path = require('node:path');

const routes = [
  ['get', '/v1/accounts', 'List accounts'],
  ['get', '/v1/accounts/{id}/entries', 'Ledger statement for an account'],
  ['post', '/v1/transfers', 'Create a transfer (requires Idempotency-Key)'],
  ['get', '/v1/invoices', 'List invoices'],
  ['get', '/v1/cards', 'List cards'],
];

const paths = {};
for (const [method, route, summary] of routes) {
  paths[route] = {
    ...paths[route],
    [method]: {
      summary,
      security: [{ bearer: [] }],
      responses: { 200: { description: 'OK' }, 401: { description: 'Unauthorized' } },
    },
  };
}

const spec = {
  openapi: '3.0.3',
  info: { title: 'FinFlow API', version: require('../package.json').version },
  components: { securitySchemes: { bearer: { type: 'http', scheme: 'bearer' } } },
  paths,
};
fs.writeFileSync(path.join(__dirname, '..', 'docs', 'openapi.json'), JSON.stringify(spec, null, 2));
console.log(`wrote ${Object.keys(paths).length} paths`);
