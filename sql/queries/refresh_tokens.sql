-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at) VALUES (
    $1,
    Now(),
    Now(),
    $2,
    Now() + INTERVAL '10 minutes',
    NULL
) RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT users.id FROM users
JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
AND refresh_tokens.revoked_at IS NULL
AND expires_at > now();

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = Now()
WHERE token = $1;