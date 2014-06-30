
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE users ADD COLUMN recovery_code VARCHAR(30) NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE users DROP COLUMN recovery_code;
