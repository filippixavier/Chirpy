-- name: MakeRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefrehToken :one
UPDATE refresh_tokens
SET revoked_at = $2, updated_at = NOW()
WHERE token = $1
RETURNING *;