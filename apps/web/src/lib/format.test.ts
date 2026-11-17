import test from 'node:test';
import assert from 'node:assert/strict';
import { formatMinor, maskIban, relativeTime } from './format';

test('formatMinor renders currency', () => {
  assert.equal(formatMinor(123456, 'USD'), '$1,234.56');
});

test('maskIban hides the middle', () => {
  assert.equal(maskIban('DE89370400440532013000'), 'DE89 •••• 3000');
  assert.equal(maskIban(undefined), '—');
});

test('relativeTime returns days for old timestamps', () => {
  const now = Date.parse('2025-03-10T00:00:00Z');
  assert.equal(relativeTime('2025-03-08T00:00:00Z', now), '2 days ago');
});
