import type { CreateTransferRequest } from '@finflow/types';

export type FieldErrors<T> = Partial<Record<keyof T, string>>;

export function validateTransfer(
  req: Partial<CreateTransferRequest>,
  availableMinor: number,
): FieldErrors<CreateTransferRequest> {
  const errors: FieldErrors<CreateTransferRequest> = {};
  if (!req.sourceAccountId) errors.sourceAccountId = 'Choose a source account';
  if (!req.destinationAccountId) errors.destinationAccountId = 'Choose a destination';
  if (req.sourceAccountId && req.sourceAccountId === req.destinationAccountId) {
    errors.destinationAccountId = 'Destination must differ from source';
  }
  if (!req.amountMinor || req.amountMinor <= 0) errors.amountMinor = 'Enter an amount greater than zero';
  else if (req.amountMinor > availableMinor) errors.amountMinor = 'Insufficient available balance';
  if (req.memo && req.memo.length > 140) errors.memo = 'Memo is limited to 140 characters';
  return errors;
}
