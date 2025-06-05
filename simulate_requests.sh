#!/bin/bash

# ---------------------------
# Configurable Parameters
# ---------------------------
URL="http://localhost:8080/check"
KEY="load_test_user"
RATE=1
INTERVAL=60       # seconds
REQUESTS=1000       # number of requests to send
DELAY_MS=0      # milliseconds between requests

# ---------------------------
# Simulation Logic
# ---------------------------
success=0
rejected=0

echo "🔁 Simulating $REQUESTS requests to $URL"

for i in $(seq 1 $REQUESTS); do
  response=$(curl -s -X POST $URL \
    -H "Content-Type: application/json" \
    -d "{\"key\": \"$KEY\", \"rate\": $RATE, \"interval\": $INTERVAL}")

  allowed=$(echo "$response" | jq -r '.allowed')
  retry_after=$(echo "$response" | jq -r '.retry_after_sec')

  if [[ "$allowed" == "true" ]]; then
    echo "[$i] ✅ Allowed"
    ((success++))
  else
    echo "[$i] ❌ Rejected (retry after ${retry_after}s)"
    ((rejected++))
  fi

  sleep $(bc <<< "scale=3; $DELAY_MS / 1000")
done

# ---------------------------
# Summary
# ---------------------------
echo "----------------------------------"
echo "✔️ Success: $success"
echo "❌ Rejected: $rejected"
echo "📊 Rate: $RATE req / $INTERVAL sec"