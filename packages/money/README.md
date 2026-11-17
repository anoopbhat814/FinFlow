# @finflow/money

Immutable money value object for TypeScript. Amounts are stored as safe integers in minor units, so there is no floating point drift when adding, subtracting or splitting.

```ts
import { Money } from '@finflow/money';

const price = Money.parse('19.99', 'USD');
const [a, b, c] = price.allocate(3); // 6.67, 6.66, 6.66
price.add(Money.fromMinor(1, 'USD')).format(); // "$20.00"
```

## Rules

- Operations across different currencies throw `TypeError`.
- `parse` rejects more decimal places than the currency exponent allows (JPY has none).
- `allocate` never loses or creates a unit: the parts always sum to the original amount.

The Go counterpart lives in `services/ledger/pkg/money` and follows the same rounding rules.

## Tests

```
npm test --workspace @finflow/money
```
