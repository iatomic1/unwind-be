-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS book (
  id UUID PRIMARY KEY,
  title VARCHAR(255) NOT NULL,  -- tags:`validate:"min=10,max=255"`
  author VARCHAR(255) NOT NULL,  -- tags:`validate:"required"`
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- Corrected function with properly terminated dollar-quoted string
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON book
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS book;
DROP FUNCTION IF EXISTS update_updated_at_column;
-- +goose StatementEnd
