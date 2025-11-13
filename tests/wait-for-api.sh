#!/bin/bash
set -euo pipefail

URL="${1:-http://api:8081/games}"
TRIES=60
SLEEP=2

echo "⏳ Esperando API en: $URL"

for i in $(seq 1 "$TRIES"); do
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$URL" || true)
  if [ "$STATUS" = "200" ]; then
    echo "✅ API lista (HTTP 200)"
    exit 0
  fi
  echo "… aún no lista (status=$STATUS). intento $i/$TRIES"
  sleep "$SLEEP"
done

echo "❌ Timeout esperando la API en $URL"
exit 1