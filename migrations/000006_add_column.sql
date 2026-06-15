-- +goose Up
ALTER TABLE certificate_applications ADD COLUMN comment VARCHAR(255) DEFAULT NULL;

-- +goose Down
ALTER TABLE certificate_applications DROP COLUMN comment;