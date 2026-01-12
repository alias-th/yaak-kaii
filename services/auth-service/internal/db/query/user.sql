-- name: CreateUser :one
INSERT INTO
    users (
        email,
        password_hash,
        first_name,
        last_name,
        phone_number
    )
VALUES
    ($1, $2, $3, $4, $5) RETURNING *;


-- name: GetUserById :one
SELECT * FROM users u WHERE u.id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users u WHERE u.email = $1 LIMIT 1;