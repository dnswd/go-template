-- name: GetUser :one
SELECT *
FROM users 
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (email, name, balance) 
VALUES ($1, $2, $3) 
RETURNING *;

-- name: DeleteUser :execrows
DELETE FROM users WHERE id = $1;
