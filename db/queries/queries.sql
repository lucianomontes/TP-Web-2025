-- name: GetGame :one
SELECT id, titulo, descripcion, categoria, to_char(fecha, 'YYYY-MM-DD') AS fecha, estado, imagen, created_at
FROM games
WHERE id = $1;

-- name: ListGames :many
SELECT id, titulo, descripcion, categoria, to_char(fecha, 'YYYY-MM-DD') AS fecha, estado, imagen, created_at
FROM games
ORDER BY titulo;

-- name: ListWantedGames :many
SELECT id, titulo, descripcion, categoria, to_char(fecha, 'YYYY-MM-DD') AS fecha, estado, imagen, created_at
FROM games
WHERE estado = 'deseado'
ORDER BY titulo;

-- name: CreateGame :one
INSERT INTO games (titulo, descripcion, categoria, fecha, estado, imagen)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, titulo, descripcion, categoria, to_char(fecha, 'YYYY-MM-DD') AS fecha, estado, imagen, created_at;

-- name: UpdateGame :one
UPDATE games
SET titulo = $2, descripcion = $3, categoria = $4, fecha = $5, estado = $6, imagen = $7
WHERE id = $1
RETURNING *;

-- name: UpdateGameState :one
UPDATE games
SET estado = $2
WHERE id = $1
RETURNING *;

-- name: DeleteGame :one
DELETE FROM games
WHERE id = $1
RETURNING *;