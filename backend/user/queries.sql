-- name: GetUser :one
SELECT *
FROM users 
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (email, name) 
VALUES ($1, $2) 
RETURNING *;

-- name: DeleteUser :execrows
DELETE FROM users WHERE id = $1;
