-- name: CreateUser :one
INSERT INTO users (
    id,
    name,
    created_at,
    updated_at,
    email,
    password_hash
)
VALUES (
    gen_random_uuid(),
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;

-- name: GetUserFromEmail :one
SELECT * FROM USERS
WHERE email = $1;