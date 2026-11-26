-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

create extension IF NOT EXISTS "pgcrypto";

CREATE TABLE TENANTS
(

    "id"         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "issuer_id"  INTEGER,
    "dialect"    VARCHAR(255),
    "uri"        VARCHAR(255),
    "user"       VARCHAR(255),
    "password"   VARCHAR(255),
    "type"       VARCHAR(255),
    "created_at" TIMESTAMP,
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP

);

CREATE INDEX idx_tenants_issuer_id ON tenants USING btree (issuer_id);
CREATE INDEX idx_tenants_deleted_at ON tenants USING btree (deleted_at);
CREATE INDEX idx_tenants_type_access ON tenants USING btree ("type");

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE TENANTS;
