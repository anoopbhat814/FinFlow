'use strict';
const test = require('node:test');
const assert = require('node:assert/strict');
const { render } = require('./templates');
const { backoff } = require('./queue');
const { sign } = require('./webhooks');

test('render formats currency amounts', () => {
  const { subject } = render('transfer_completed', { amountMinor: 12550, currency: 'USD', recipient: 'Bob' });
  assert.equal(subject, 'You sent $125.50');
});

test('render rejects unknown templates', () => {
  assert.throws(() => render('nope', {}), /unknown template/);
});

test('backoff grows quadratically', () => {
  assert.deepEqual([1, 2, 3].map(backoff), [1000, 4000, 9000]);
});

test('webhook signatures are deterministic', () => {
  assert.equal(sign('s', 1, '{}'), sign('s', 1, '{}'));
  assert.notEqual(sign('s', 1, '{}'), sign('s', 2, '{}'));
});
