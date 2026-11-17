import test from 'node:test';
import assert from 'node:assert/strict';
import type { Transfer } from '@finflow/types';
import { transfersToCsv } from './csv';

test('transfersToCsv escapes commas and quotes in memos', () => {
  const t = { id: '1', createdAt: '2025-01-01', status: 'COMPLETED', currency: 'USD', amountMinor: 1050, memo: 'Lunch, "team"' } as Transfer;
  const lines = transfersToCsv([t]).split('\n');
  assert.equal(lines[1], '1,2025-01-01,COMPLETED,USD,10.50,"Lunch, ""team"""');
});
