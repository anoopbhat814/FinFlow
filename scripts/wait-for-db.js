#!/usr/bin/env node
'use strict';
const net = require('node:net');

const url = new URL(process.env.DATABASE_URL || 'postgresql://localhost:5432/finflow');
const deadline = Date.now() + Number(process.argv[2] || 30) * 1000;

function attempt() {
  const sock = net.connect(Number(url.port || 5432), url.hostname);
  sock.once('connect', () => { sock.end(); console.log('database is accepting connections'); });
  sock.once('error', () => {
    sock.destroy();
    if (Date.now() > deadline) {
      console.error('timed out waiting for database');
      process.exit(1);
    }
    setTimeout(attempt, 500);
  });
}
attempt();
