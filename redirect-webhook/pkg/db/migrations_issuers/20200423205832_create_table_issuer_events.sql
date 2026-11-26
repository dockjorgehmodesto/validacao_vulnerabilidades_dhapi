-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE ISSUER_EVENTS
(

    "id"            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name"          VARCHAR(255) NOT NULL,
    "connection_id" UUID NOT NULL,
    "issuer_id"     INT          NOT NULL,
    "created_at"    TIMESTAMP,
    "updated_at"    TIMESTAMP,
    "deleted_at"    TIMESTAMP,

    FOREIGN KEY (connection_id) REFERENCES AMQP_CONNECTIONS(id)

);

CREATE INDEX idx_issuer_events_deleted_at ON ISSUER_EVENTS USING btree (deleted_at);
CREATE INDEX idx_issuer_events_issuer_id ON ISSUER_EVENTS USING btree (issuer_id);
CREATE INDEX idx_issuer_events_issuer_id_name ON ISSUER_EVENTS USING btree (issuer_id, name);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE ISSUER_EVENTS;

