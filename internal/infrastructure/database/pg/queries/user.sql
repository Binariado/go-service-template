-- name: GetUserByID :one
SELECT id, name, created_at, updated_at
FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (name)
VALUES ($1)
RETURNING id, name, created_at, updated_at;
