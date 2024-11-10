-- name: InsertProfile :one
INSERT INTO profile (
  id, username, user_id
) VALUES ( uuid_generate_v4(), $1, $2 )
RETURNING *;


-- name: UpdateProfile :one
UPDATE profile
SET 
    username = COALESCE(sqlc.arg('username'), username),
    name = COALESCE(sqlc.narg('name'), name),
    cover_pic = COALESCE(sqlc.narg('cover_pic'), cover_pic),
    profile_pic = COALESCE(sqlc.narg('profile_pic'), profile_pic),
    updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING *;


-- name: GetProfileByUserId :one
SELECT * FROM profile WHERE user_id = $1;
