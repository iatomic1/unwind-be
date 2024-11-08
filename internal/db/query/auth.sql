-- name: RegisterUser :one 
INSERT INTO "user" (
 id, email, password_hash 
) VALUES ( uuid_generate_v4(), $1, $2 )
RETURNING *
;
