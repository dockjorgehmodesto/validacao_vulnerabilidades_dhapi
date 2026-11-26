-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE web_hooks
    ALTER COLUMN url TYPE text;

ALTER TABLE web_hook_headers
    ALTER COLUMN value TYPE text;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
