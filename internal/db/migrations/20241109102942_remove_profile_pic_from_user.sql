-- +goose Up
-- +goose StatementBegin
ALTER TABLE "user"
DROP COLUMN profile_pic;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user"
ADD COLUMN profile_pic TEXT;
-- +goose StatementEnd
