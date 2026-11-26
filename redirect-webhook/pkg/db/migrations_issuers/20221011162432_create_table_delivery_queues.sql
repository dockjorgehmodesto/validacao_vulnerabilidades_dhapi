-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE DELIVERY_QUEUES
(

    "id"                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "issuer_id"         INT NOT NULL,
    "event_type"        VARCHAR(255),
    "exchange_name"     VARCHAR(255),
    "queue_name"        varchar(255),
    "queue_key"         varchar(255),
    "enabled"           BOOLEAN,
    "created_at"        TIMESTAMP,
    "updated_at"        TIMESTAMP,
    "deleted_at"        TIMESTAMP

);

CREATE INDEX idx_deliveryqueue_issuer_id ON DELIVERY_QUEUES USING btree (issuer_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE DELIVERY_QUEUES;