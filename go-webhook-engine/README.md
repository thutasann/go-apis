# Go Webhook Engine

## Test Webhook

```bash
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "idempotency_key": "abc-123",
    "type": "user.created",
    "payload": { "user_id": 42 }
  }'
```
