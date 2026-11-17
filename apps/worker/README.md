# @finflow/worker

Background consumer for the `notifications` Redis stream. It turns domain events emitted by the ledger service into emails and signed webhook deliveries.

## Running locally

```
cp .env.example .env
npm run start --workspace @finflow/worker
```

## Job shape

```json
{
  "channel": "EMAIL",
  "to": "ada@example.com",
  "template": "transfer_completed",
  "data": { "amountMinor": 12500, "currency": "USD", "recipient": "Grace" }
}
```

## Retry policy

Failed jobs are re-queued with quadratic backoff (1s, 4s, 9s, 16s). After five attempts the payload is moved to `notifications:dead` for manual inspection.

## Webhook signatures

Every webhook carries `FinFlow-Timestamp` and `FinFlow-Signature: v1=<hex>`, where the signature is `HMAC-SHA256(secret, "<timestamp>.<body>")`. Receivers should reject timestamps older than five minutes.
