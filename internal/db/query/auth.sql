-- name: RegisterUser :one 
INSERT INTO "user" (
 id, email, password 
) VALUES ( uuid_generate_v4(), $1, $2 )
RETURNING *
;

-- name: GetUserByEmail :one
SELECT * FROM "user"
WHERE email = $1;
