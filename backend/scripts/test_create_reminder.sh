#!/bin/bash

curl -s -X POST http://localhost:8080/webhook \
-H "Content-Type: application/json" \
-d '{
  "from":"628987654321",
  "body":"remind me to test scheduler now"
}'

echo ""