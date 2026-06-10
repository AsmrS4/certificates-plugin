-- +goose Up
ALTER TABLE certificates ADD COLUMN file_id TEXT DEFAULT NULL;

-- +goose Down
ALTER TABLE certificates DROP COLUMN file_id;