# API reference (v1)

All endpoints require `Authorization: Bearer <jwt>`.

| Method | Path | Description |
|--------|------|-------------|
| GET | `/v1/accounts` | List the caller's accounts |
| GET | `/v1/accounts/{id}/entries?limit=50` | Ledger statement |
| POST | `/v1/transfers` | Create a transfer; `Idempotency-Key` header required |
| POST | `/v1/transfers/{id}/refunds` | Refund part or all of a completed transfer (`amountMinor` 0 = remaining) |
| GET | `/v1/accounts/{id}/statements/{year}/{month}` | Monthly statement for a closed month |
| GET | `/v1/cards` | List the caller's cards |
| POST | `/v1/cards/{id}/freeze`, `/unfreeze` | Freeze or unfreeze an active card |
| POST | `/v1/cards/{id}/authorizations` | Evaluate a purchase against status, expiry and monthly limit |
| GET | `/v1/invoices?status=` | List invoices with computed totals |
| GET | `/v1/invoices/{id}` | Fetch one invoice |
| POST | `/v1/invoices/{id}/pay` | Mark a sent or overdue invoice as paid |
| GET, POST | `/v1/webhooks` | List or register webhook endpoints (secret returned once on create) |

## Errors

Errors are JSON: `{ "code": "unprocessable", "message": "insufficient funds" }`. Every response carries an `X-Request-ID`. Rate-limited calls return `429` with `Retry-After`.

## Example

```bash
curl -X POST $API/v1/transfers \
  -H "Authorization: Bearer $TOKEN" -H "Idempotency-Key: $(uuidgen)" \
  -d '{"sourceAccountId":"a1","destinationAccountId":"a2","amountMinor":1250,"currency":"USD"}'
```
