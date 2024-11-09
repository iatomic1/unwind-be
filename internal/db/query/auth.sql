-- name: RegisterUser :one 
INSERT INTO "user" (
 id, email, password 
) VALUES ( uuid_generate_v4(), $1, $2 )
RETURNING *
;

-- name: GetUserByEmail :one
SELECT * FROM "user"
WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM "user"
WHERE id = $1;

