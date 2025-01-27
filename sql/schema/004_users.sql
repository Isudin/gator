-- +goose Up
ALTER TABLE users
ADD CONSTRAINT username UNIQUE (name);

-- +goose Down
ALTER TABLE users
DROP CONSTRAINT username;