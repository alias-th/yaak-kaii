-- name: CreateGuest :one
INSERT INTO
    guests (token_hash, ip_addr, user_agent, metadata)
VALUES
    ($1, $2, $3, $4) RETURNING *;