-- name: GetRefreshTokenByTokenHash :one
SELECT * FROM refresh_tokens WHERE token_hash = $1 limit 1;