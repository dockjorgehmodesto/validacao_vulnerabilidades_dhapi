-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE AMQP_CONNECTIONS ADD COLUMN PROCESSOR_TYPE VARCHAR(255);

UPDATE AMQP_CONNECTIONS SET PROCESSOR_TYPE = 'CDC_MESSAGE_PROCESSOR';

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE AMQP_CONNECTIONS DROP COLUMN PROCESSOR_TYPE;
