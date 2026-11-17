# API reference (v1)

All endpoints require `Authorization: Bearer <jwt>`.

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/accounts` | List the caller's accounts |
| GET | `/v1/accounts/{id}/entries?limit=50` | Ledger statement |
| POST | `/v1/transfers` | Create a transfer; `Idempotency-Key` header required |

## Errors

Errors are JSON: `{ "code": "unprocessable", "message": "insufficient funds" }`. Rate-limited calls return `429` with `Retry-After`.

## Example

```bash
curl -X POST $API/v1/transfers \
  -H "Authorization: Bearer $TOKEN" -H "Idempotency-Key: $(uuidgen)" \
  -d '{"sourceAccountId":"a1","destinationAccountId":"a2","amountMinor":1250,"currency":"USD"}'
```
