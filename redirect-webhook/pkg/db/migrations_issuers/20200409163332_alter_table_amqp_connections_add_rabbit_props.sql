-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE amqp_connections ADD COLUMN IF NOT EXISTS "broker" VARCHAR(255);
ALTER TABLE amqp_connections ADD COLUMN IF NOT EXISTS "virtual_host" VARCHAR(255);
ALTER TABLE amqp_connections ADD COLUMN IF NOT EXISTS "exchange_name" VARCHAR(255);
ALTER TABLE amqp_connections ADD COLUMN IF NOT EXISTS "queue_name" VARCHAR(255);
ALTER TABLE amqp_connections ADD COLUMN IF NOT EXISTS "queue_key" VARCHAR(255);

UPDATE amqp_connections
SET broker = 'SERVICEBUS' WHERE amqp_connections.broker IS NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE amqp_connections DROP COLUMN IF EXISTS "broker";
ALTER TABLE amqp_connections DROP COLUMN IF EXISTS "virtual_host";
ALTER TABLE amqp_connections DROP COLUMN IF EXISTS "exchange_name";
ALTER TABLE amqp_connections DROP COLUMN IF EXISTS "queue_name";
ALTER TABLE amqp_connections DROP COLUMN IF EXISTS "queue_key";