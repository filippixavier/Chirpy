-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE users.id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE users.email = $1;

-- name: ClearUsers :exec
DELETE FROM users;

-- name: UpdateUser :one
UPDATE users
SET email = $2, hashed_password = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpgradeRedStatus :one
UPDATE users
SET is_chirpy_red = $2
WHERE id = $1
RETURNING *;