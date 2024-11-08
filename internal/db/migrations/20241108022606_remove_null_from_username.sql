-- +goose Up
-- +goose StatementBegin
ALTER TABLE "user"
    ALTER COLUMN username DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user"
    ALTER COLUMN username SET NOT NULL;
-- +goose StatementEnd
