-- name: GetUserWatchList :many
SELECT * FROM watch_list
WHERE user_id = $1;


-- name: GetWatchListByMediaID :one
SELECT * FROM watch_list
WHERE media_id = $1 AND user_id = $2;

-- name: AddToList :one
INSERT INTO watch_list (
  id, user_id, type, media_id, poster, title, status, episodes, duration, media_type 
) VALUES ( uuid_generate_v4(), $1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateWatchListStatus :one
UPDATE watch_list
SET
    status = COALESCE(sqlc.arg('status'), status),
    updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteWatchList :one
DELETE FROM watch_list
  WHERE id = $1
RETURNING *;
