-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE WEB_HOOKS
(
    "id"         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "event_type" VARCHAR(255),
    "url"        VARCHAR(255),
    "enabled"    BOOLEAN,
    "issuer_id"  INT NOT NULL,
    "created_at" TIMESTAMP,
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP
);

CREATE TABLE WEB_HOOK_HEADERS
(
    "id"          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "web_hook_id" UUID NOT NULL,
    "key"         VARCHAR(255),
    "value"       VARCHAR(255),
    "created_at"  TIMESTAMP,
    "updated_at"  TIMESTAMP,
    "deleted_at"  TIMESTAMP
);

CREATE INDEX idx_webhook_type ON WEB_HOOKS USING btree (event_type);
CREATE INDEX idx_webhook_deleted_at ON WEB_HOOKS USING btree (deleted_at);
CREATE INDEX idx_wehook_issuer_id ON WEB_HOOKS USING btree (issuer_id);
CREATE INDEX idx_wehook_issuer_id_enabled ON WEB_HOOKS USING btree (issuer_id, enabled);
CREATE INDEX idx_webhookheader_deleted_at ON WEB_HOOK_HEADERS USING btree (deleted_at);
CREATE INDEX idx_webhookheader_webhook ON WEB_HOOK_HEADERS USING btree (web_hook_id);
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE WEB_HOOK_HEADERS;
DROP TABLE WEB_HOOKS;