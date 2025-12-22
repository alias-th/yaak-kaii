-- name: CreateGuest :one
INSERT INTO
    guests (token_hash, ip_addr, user_agent, expires_at, metadata)
VALUES
    ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetGuestByToken :one
SELECT * FROM guests g WHERE g.token_hash = $1 limit 1;