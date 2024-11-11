-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "watch_list" (
  id UUID PRIMARY KEY,
  user_id UUID UNIQUE NOT NULL,
  anilist_id VARCHAR(16),  -- tags:`binding:"required"`
  hianime_id VARCHAR(256),  -- tags:`binding:"required"`
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT anilist_id_unique UNIQUE (anilist_id),
  CONSTRAINT hi_anime_id_unique UNIQUE (hianime_id),
  CONSTRAINT fk_user
    FOREIGN KEY(user_id)
    REFERENCES "user"(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "watch_list";
-- +goose StatementEnd
