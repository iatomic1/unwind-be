-- name: GetUserWatchList :many
SELECT * FROM watch_list
WHERE user_id = $1;

-- name: AddToList :one
INSERT INTO watch_list (
  id, anilist_id, hianime_id, user_id
) VALUES ( uuid_generate_v4(), $1, $2, $3 )
RETURNING *;
