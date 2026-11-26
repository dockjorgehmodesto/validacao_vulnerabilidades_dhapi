-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

create extension IF NOT EXISTS "pgcrypto";

CREATE TABLE KEY_VAULTS
(

    "id"                         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "issuer_public_key"          bytea,
    "aes_key"                    bytea,
    "issuer_key_expiration_date" TIMESTAMP,
    "aes_key_expiration_date"    TIMESTAMP,
    "enabled"                    BOOLEAN,
    "issuer_id"                  INT NOT NULL,
    "created_at"                 TIMESTAMP,
    "updated_at"                 TIMESTAMP,
    "deleted_at"                 TIMESTAMP

);

CREATE INDEX idx_keyvault_deleted_at ON KEY_VAULTS USING btree (deleted_at);
CREATE INDEX idx_keyvault_issuer_id ON KEY_VAULTS USING btree (issuer_id);
CREATE INDEX idx_keyvault_issuer_id_enabled ON KEY_VAULTS USING btree (issuer_id, enabled);
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE KEY_VAULTS;