-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS profile (
  id UUID PRIMARY KEY,
  user_id UUID UNIQUE NOT NULL, -- tags:`binding:"required,uuid"`
  profile_pic TEXT,
  name VARCHAR(100),
  username VARCHAR(20) NOT NULL, -- tags:`binding:"required,min=8"`
  cover_pic TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_user
    FOREIGN KEY(user_id)
    REFERENCES "user"(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT profile_username_unique UNIQUE (username)
);

CREATE INDEX idx_profile_user_id ON profile(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- ALTER TABLE profile
DROP TABLE IF EXISTS profile;
-- +goose StatementEnd
