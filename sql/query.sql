-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserWithEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY email;

-- name: CreateUser :one
INSERT INTO users(name, email, is_admin)
VALUES($1, $2, $3)
RETURNING *;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;

-- name: MakeAdmin :exec
UPDATE users 
SET is_admin = TRUE
WHERE email = $1;

-- name: CreateReceipt :one
INSERT INTO receipts(user_id, fpd, total, time, optype, place)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListReceipts :many
SELECT * 
FROM receipts
WHERE user_id = $1
ORDER BY time ;
