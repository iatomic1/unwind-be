-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS book ( id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  author VARCHAR(255),
  plublished_date DATE );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS book;
-- +goose StatementEnd
