-- +goose Up
ALTER TABLE certificates RENAME COLUMN receiver_id TO order_id;

-- +goose Down
ALTER TABLE certificates RENAME COLUMN order_id TO receiver_id;