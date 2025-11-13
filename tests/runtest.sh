#!/bin/bash
set -euo pipefail

API_URL="http://api:8081"

request() {
  local method="$1"
  local url="$2"
  local title="$3"
  local data="${4-}"

  echo -e "\n--- $title ---"

  if [ -z "$data" ]; then
    STATUS=$(curl -s -o /tmp/body.json -w "%{http_code}" -X "$method" "$url")
  else
    STATUS=$(curl -s -o /tmp/body.json -w "%{http_code}" \
      -X "$method" "$url" \
      -H "Content-Type: application/json" \
      -d "$data")
  fi

  echo "Status: $STATUS"
  LAST_STATUS="$STATUS" # ← global para asserts

  if jq -e . >/dev/null 2>&1 < /tmp/body.json; then
    jq . < /tmp/body.json
  else
    cat /tmp/body.json
  fi
}

expect_status() {
  local expected="$1"
  local label="$2"
  if [[ "$LAST_STATUS" != "$expected" ]]; then
    echo "❌ FAIL $label — esperado $expected pero fue $LAST_STATUS"
    exit 1
  else
    echo "✅ OK $label"
  fi
}

# GET /games (inicio)
request "GET" "$API_URL/games" "Test GET /games (inicio)"
expect_status "200" "GET /games (inicio)"

# POST /games (válido)
POST_DATA='{
  "titulo": "Mario Bros",
  "descripcion": "Aventura y plataformas",
  "categoria": "Plataformas",
  "fecha": "2024-10-30",
  "estado": "none",
  "imagen": "img/mario.png"
}'
request "POST" "$API_URL/games" "Test POST /games (válido)" "$POST_DATA"
expect_status "201" "POST /games (válido)"

# Extraer ID desde el último body
GAME_ID=$(jq -r '.id // empty' < /tmp/body.json)
echo -e "\n📌 ID creado para pruebas: $GAME_ID\n"
if [ -z "$GAME_ID" ] || [ "$GAME_ID" = "null" ]; then
  echo "❌ ERROR: No se pudo obtener GAME_ID del POST"
  exit 1
fi

# GET /games (después del POST)
request "GET" "$API_URL/games" "Test GET /games (después del POST)"
expect_status "200" "GET /games (después del POST)"

# GET /games/{id}
request "GET" "$API_URL/games/$GAME_ID" "Test GET /games/$GAME_ID"
expect_status "200" "GET /games/{id}"

# PUT /games/{id} (válido)
PUT_DATA_VALID='{
  "titulo": "Mario Bros Deluxe",
  "descripcion": "Mejorado y remasterizado",
  "categoria": "Plataformas",
  "fecha": "2024-10-31",
  "estado": "none",
  "imagen": "img/mario_rem.png"
}'
request "PUT" "$API_URL/games/$GAME_ID" "Test PUT /games/$GAME_ID (válido)" "$PUT_DATA_VALID"
expect_status "200" "PUT /games/{id} (válido)"

# PUT /game_state/{id}?state=deseado (válido)
request "PUT" "$API_URL/game_state/$GAME_ID?state=deseado" "Test PUT state=deseado"
expect_status "200" "PUT /game_state (válido)"

# TEST QUE ESPERAN ERROR

# ❌ POST /games con campo vacío (titulo = "")
POST_EMPTY_FIELD='{
  "titulo": "",
  "descripcion": "desc ok",
  "categoria": "Plataformas",
  "fecha": "2024-10-30",
  "estado": "none",
  "imagen": "img/imagen_ok.png"
}'
request "POST" "$API_URL/games" "NEG: POST con campo vacío (titulo)" "$POST_EMPTY_FIELD"
expect_status "400" "NEG POST campo vacío"

# ❌ POST /games con estado inválido
POST_BAD_STATE='{
  "titulo": "Juego Estado Invalido",
  "descripcion": "desc",
  "categoria": "Aventura",
  "fecha": "2024-11-01",
  "estado": "invalido",
  "imagen": "img/juego.png"
}'
request "POST" "$API_URL/games" "NEG: POST con estado inválido" "$POST_BAD_STATE"
expect_status "400" "NEG POST estado inválido"

# ❌ PUT /games/{id} con estado inválido
PUT_BAD_STATE='{
  "titulo": "Mario Bros Deluxe",
  "descripcion": "desc",
  "categoria": "Plataformas",
  "fecha": "2024-10-31",
  "estado": "malo",
  "imagen": "img/mario_rem.png"
}'
request "PUT" "$API_URL/games/$GAME_ID" "NEG: PUT /games/{id} estado inválido" "$PUT_BAD_STATE"
expect_status "500" "NEG PUT /games/{id} estado inválido"

# ❌ PUT /game_state/{id}?state=invalido
request "PUT" "$API_URL/game_state/$GAME_ID?state=invalido" "NEG: PUT /game_state estado inválido"
expect_status "500" "NEG PUT /game_state estado inválido"

#####################################################

# GET /wanted_games
request "GET" "$API_URL/wanted_games" "Test GET /wanted_games"
expect_status "200" "GET /wanted_games"

# DELETE /games/{id}
request "DELETE" "$API_URL/games/$GAME_ID" "Test DELETE /games/$GAME_ID"
expect_status "204" "DELETE /games/{id}"

echo -e "\n--- TEST COMPLETADO ✅ ---\n"