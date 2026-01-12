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

-- name: GetShopBySellerID :one
SELECT
    *
FROM
    shops
WHERE
    seller_id = $1;


-- name: GetShopByUserID :one
SELECT 
    s.*
FROM
    shops s
LEFT JOIN
    sellers se ON s.seller_id = se.id
WHERE
    se.user_id = $1;