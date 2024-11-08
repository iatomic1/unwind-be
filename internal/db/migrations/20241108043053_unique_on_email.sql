-- +goose Up
-- +goose StatementBegin
ALTER TABLE "user"
ADD CONSTRAINT user_email_unique UNIQUE (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user"
DROP CONSTRAINT IF EXISTS user_email_unique;
-- +goose StatementEnd
