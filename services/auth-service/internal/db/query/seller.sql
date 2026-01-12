-- name: CreateSeller :one
INSERT INTO
    sellers (
     user_id,
     name
    )
VALUES
    ($1, $2) RETURNING *;

-- name: GetSellerByUserID :one
SELECT
    *
FROM
    sellers
WHERE
    user_id = $1;