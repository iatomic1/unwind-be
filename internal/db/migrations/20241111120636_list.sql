-- +goose Up
-- +goose StatementBegin
CREATE TYPE valid_types AS ENUM ('anime', 'movie', 'kdrama', 'manga');
CREATE TYPE status AS ENUM ('watching', 'on-hold', 'planning', 'dropped', 'completed');
CREATE TYPE media_type AS ENUM ('tv', 'movie');
CREATE TABLE IF NOT EXISTS "watch_list" (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  type VALID_TYPES NOT NULL,  -- tags:`binding:"required"`
  media_type MEDIA_TYPE NOT NULL, -- tags:`binding:"required"`
  media_id VARCHAR(256),  -- tags:`binding:"required"`
  poster TEXT NOT NULL,  -- tags:`binding:"required"`
  title VARCHAR(256) NOT NULL,  -- tags:`binding:"required"`
  status STATUS NOT NULL,  -- tags:`binding:"required"`
  rated INT,
  episodes INT,
  duration INT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_user
    FOREIGN KEY(user_id)
    REFERENCES "user"(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
CREATE INDEX idx_watch_list_user_id ON "watch_list" (user_id);
CREATE INDEX idx_watch_list_type ON "watch_list" (type);
CREATE INDEX idx_watch_list_status ON "watch_list" (status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "watch_list";
DROP TYPE IF EXISTS valid_types;
DROP TYPE IF EXISTS status;
DROP TYPE IF EXISTS media_type;
-- +goose StatementEnd
