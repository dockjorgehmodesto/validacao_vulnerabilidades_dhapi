-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE WEB_HOOKS DROP COLUMN ENABLED;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE WEB_HOOKS ADD COLUMN ENABLED BOOL;
