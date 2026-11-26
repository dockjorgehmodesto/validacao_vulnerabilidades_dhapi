-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE AMQP_CONNECTIONS
(

    "id"                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "connection_string" VARCHAR(255),
    "topic"             VARCHAR(255),
    "subscription_name" varchar(255),
    "enabled"           BOOLEAN,
    "issuer_id"         INT NOT NULL,
    "created_at"        TIMESTAMP,
    "updated_at"        TIMESTAMP,
    "deleted_at"        TIMESTAMP

);

CREATE INDEX idx_amqpconnection_deleted_at ON AMQP_CONNECTIONS USING btree (deleted_at);
CREATE INDEX idx_amqpconnection_issuer_id ON AMQP_CONNECTIONS USING btree (issuer_id);
CREATE INDEX idx_amqpconnection_issuer_id_enabled ON AMQP_CONNECTIONS USING btree (issuer_id, enabled);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE AMQP_CONNECTIONS;

