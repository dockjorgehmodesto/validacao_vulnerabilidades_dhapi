-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE issuer_events ALTER COLUMN "connection_id" DROP NOT NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE issuer_events ALTER COLUMN "connection_id" SET NOT NULL;