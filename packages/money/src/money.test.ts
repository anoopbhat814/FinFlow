import test from 'node:test';
import assert from 'node:assert/strict';
import { Money } from './index';

test('parse handles decimals and zero-exponent currencies', () => {
  assert.equal(Money.parse('12.5', 'USD').minor, 1250);
  assert.equal(Money.parse('500', 'JPY').minor, 500);
});

test('parse rejects excess precision', () => {
  assert.throws(() => Money.parse('1.234', 'USD'), RangeError);
});

test('allocate distributes remainder without losing cents', () => {
  const parts = Money.fromMinor(100, 'USD').allocate(3).map((m) => m.minor);
  assert.deepEqual(parts, [34, 33, 33]);
  assert.equal(parts.reduce((a, b) => a + b, 0), 100);
});

test('add rejects mixed currencies', () => {
  assert.throws(() => Money.fromMinor(1, 'USD').add(Money.fromMinor(1, 'EUR')), TypeError);
});
