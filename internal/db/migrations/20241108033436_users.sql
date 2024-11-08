-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
  id UUID PRIMARY KEY,
  name TEXT,
  username VARCHAR(16),
  email TEXT NOT NULL,    -- tags:`binding:"required,email"`
  password TEXT NOT NULL,    -- tags:`binding:"required"`
  profile_pic TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "user";
-- +goose StatementEnd
