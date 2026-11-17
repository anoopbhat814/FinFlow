#!/usr/bin/env node
'use strict';
// Verifies that every variable listed in .env.example is set (in .env or the shell).
const fs = require('node:fs');
const path = require('node:path');

const root = path.resolve(__dirname, '..');
const parse = (file) =>
  Object.fromEntries(
    fs.readFileSync(file, 'utf8').split('\n')
      .filter((l) => l.trim() && !l.startsWith('#'))
      .map((l) => [l.slice(0, l.indexOf('=')), l.slice(l.indexOf('=') + 1)]),
  );

const required = Object.keys(parse(path.join(root, '.env.example')));
const local = fs.existsSync(path.join(root, '.env')) ? parse(path.join(root, '.env')) : {};
const missing = required.filter((k) => !(k in local) && !process.env[k]);

if (missing.length) {
  console.error(`Missing environment variables:\n  ${missing.join('\n  ')}`);
  process.exit(1);
}
console.log(`OK: ${required.length} variables present`);
