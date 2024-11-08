-- +goose Up
-- +goose StatementBegin
ALTER TABLE "user"
   ADD COLUMN password_hash TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user"
   DROP COLUMN password_hash;
-- +goose StatementEnd
