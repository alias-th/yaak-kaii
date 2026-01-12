-- name: CreateShop :one
INSERT INTO
    shops (
     seller_id,
     name,
     slug,
     description,
     is_active
    )
VALUES
    ($1, $2, $3, $4, $5) RETURNING *;