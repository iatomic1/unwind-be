-- +goose Up
-- +goose StatementBegin
ALTER TABLE "user"
ALTER COLUMN password_hash SET NOT NULL;  -- tags:`binding:"required"

ALTER TABLE "user"
   ALTER COLUMN email SET NOT NULL;  -- tags:`binding:"required"`
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user"
   ALTER COLUMN password_hash SET NULL;

ALTER TABLE "user"
   ALTER COLUMN email SET NULL;
-- +goose StatementEnd
