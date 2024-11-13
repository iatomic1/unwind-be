-- name: InsertProfile :one
INSERT INTO profile (
  id,
  username,
  user_id,
  profile_pic,
  name,
  cover_pic
) VALUES (
  uuid_generate_v4(),
  $1,               -- username (required)
  $2,               -- user_id (required)
  $3,               -- profile_pic (optional)
  $4,               -- name (optional)
  $5                -- cover_pic (optional)
)
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

-- name: GetProfileById :one
SELECT * FROM profile WHERE id = $1;
