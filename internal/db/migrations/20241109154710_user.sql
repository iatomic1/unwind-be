-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
  id UUID PRIMARY KEY,
  email VARCHAR(100) NOT NULL,    -- tags:`binding:"required,email"`
  password VARCHAR(255) NOT NULL,    -- tags:`binding:"required"`
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT user_email_unique UNIQUE (email)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- ALTER TABLE "user"
DROP TABLE IF EXISTS "user";
-- +goose StatementEnd
