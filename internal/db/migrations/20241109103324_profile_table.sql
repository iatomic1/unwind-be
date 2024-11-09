-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS profile (
  id UUID PRIMARY KEY,
  user_id UUID UNIQUE REFERENCES "user"(id), -- tags:`binding:"required,uuid"`
  profile_pic TEXT,
  name TEXT,
  cover_pic TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS profile;
-- +goose StatementEnd
