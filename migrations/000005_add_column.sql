-- +goose Up
ALTER TABLE certificate_applications ADD COLUMN full_name TEXT DEFAULT NULL;

-- +goose Down
ALTER TABLE certificate_applications DROP COLUMN full_name;