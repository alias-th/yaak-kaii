-- name: GetRoleByName :one
SELECT * from roles r where r.name = $1 limit 1;
