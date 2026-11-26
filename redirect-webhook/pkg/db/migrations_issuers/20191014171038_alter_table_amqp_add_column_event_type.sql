-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE AMQP_CONNECTIONS ADD COLUMN EVENT_TYPE VARCHAR(255);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE AMQP_CONNECTIONS DROP COLUMN EVENT_TYPE;
