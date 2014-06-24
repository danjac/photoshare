
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE UNIQUE INDEX idx_users_upper_name ON users (UPPER(name));
CREATE INDEX idx_photos_upper_title ON photos (UPPER(title));
CREATE UNIQUE INDEX idx_tags_name ON tags (name);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP INDEX idx_users_upper_name;
DROP INDEX idx_photos_upper_title;
DROP INDEX idx_tags_name;
