-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUsers :many
SELECT id, name
FROM users
WHERE name <> '_g_invalid';

-- name: GetUserByName :one
SELECT id, name 
FROM users
WHERE name = $1;

-- name: ResetUsers :exec
DELETE FROM users
WHERE name <> '_g_invalid';
