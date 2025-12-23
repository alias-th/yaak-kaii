-- name: GetRefreshTokenByTokenHash :one
SELECT * FROM refresh_tokens WHERE token_hash = $1 limit 1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, expires_at) 
VALUES ($1, $2, $3) RETURNING *;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW() 
WHERE id = $1;

-- name: InvalidateUserTokens :exec
UPDATE refresh_tokens
SET revoked_at = NOW() 
WHERE user_id = $1 AND revoked_at IS NULL;