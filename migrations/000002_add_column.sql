-- +goose Up
ALTER TABLE certificate_applications ADD COLUMN rejection_reason TEXT DEFAULT NULL;

-- +goose Down
ALTER TABLE certificate_applications DROP COLUMN rejection_reason;