-- name: FindAllBooks :many
SELECT * FROM book;

-- name: InsertBook :one
INSERT INTO book (id, title, author)
VALUES ( uuid_generate_v4(), $1, $2 )
RETURNING *;
