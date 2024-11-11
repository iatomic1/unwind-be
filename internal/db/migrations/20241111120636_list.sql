-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "watch_list" (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  anilist_id VARCHAR(16),  -- tags:`binding:"required"`
  hianime_id VARCHAR(256),  -- tags:`binding:"required"`
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_user
    FOREIGN KEY(user_id)
    REFERENCES "user"(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
CREATE INDEX idx_watch_list_user_id ON "watch_list" (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "watch_list";
-- +goose StatementEnd
