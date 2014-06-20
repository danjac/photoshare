
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE users ADD COLUMN votes int[] DEFAULT '{}';
ALTER TABLE photos
	ADD COLUMN up_votes int DEFAULT 0,
	ADD COLUMN down_votes int DEFAULT 0;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE users DROP COLUMN votes;
ALTER TABLE photos
	DROP COLUMN up_votes,
	DROP COLUMN down_votes;

