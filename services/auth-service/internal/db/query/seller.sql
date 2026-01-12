-- name: CreateSeller :one
INSERT INTO
    sellers (
     user_id,
     name
    )
VALUES
    ($1, $2) RETURNING *;