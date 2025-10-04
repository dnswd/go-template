-- name: GetUser :one
SELECT id, email, name, created_at 
FROM users 
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (id, email, name) 
VALUES ($1, $2, $3) 
RETURNING id, email, name, created_at;

-- name: DeleteUser :execrows
DELETE FROM users WHERE id = $1;
