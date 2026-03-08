#!/bin/bash

# Use same value as WEBHOOK_SECRET in .env
curl -s -X POST http://localhost:8080/webhook \
-H "Content-Type: application/json" \
-H "X-Webhook-Secret: ${WEBHOOK_SECRET:-your_secret_here}" \
-d '{
  "from":"628987654321",
  "body":"remind me to test scheduler now"
}'

echo ""