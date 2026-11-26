-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'purchase_approved', 'DATAHUB_C6BANK_HIGH', 'pier.hooks.DATAHUB_C6BANK_HIGH', 'pier.hooks.DATAHUB_C6BANK_HIGH_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'purchase_cleared', 'DATAHUB_C6BANK_HIGH', 'pier.hooks.DATAHUB_C6BANK_HIGH', 'pier.hooks.DATAHUB_C6BANK_HIGH_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'limit', 'DATAHUB_C6BANK_HIGH', 'pier.hooks.DATAHUB_C6BANK_HIGH', 'pier.hooks.DATAHUB_C6BANK_HIGH_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'batch', 'DATAHUB_C6BANK_HIGH', 'pier.hooks.DATAHUB_C6BANK_HIGH', 'pier.hooks.DATAHUB_C6BANK_HIGH_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'purchase', 'DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'purchase_canceled', 'DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'purchase_processed', 'DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'declined', 'DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'standin', 'DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'address', 'DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'creditanalysis', 'DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'person', 'DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'phone', 'DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'adjustment', 'DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'payment', 'DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'statement', 'DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'transaction', 'DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'bankslip', 'DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'account', 'DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'card', 'DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 116, 'tokenized_purchase', 'DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR', 'pier.hooks.DATAHUB_C6BANK_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'purchase_approved', 'DATAHUB_XPI_HIGH', 'pier.hooks.DATAHUB_XPI_HIGH', 'pier.hooks.DATAHUB_XPI_HIGH_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'purchase_cleared', 'DATAHUB_XPI_HIGH', 'pier.hooks.DATAHUB_XPI_HIGH', 'pier.hooks.DATAHUB_XPI_HIGH_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'limit', 'DATAHUB_XPI_HIGH', 'pier.hooks.DATAHUB_XPI_HIGH', 'pier.hooks.DATAHUB_XPI_HIGH_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'batch', 'DATAHUB_XPI_HIGH', 'pier.hooks.DATAHUB_XPI_HIGH', 'pier.hooks.DATAHUB_XPI_HIGH_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'purchase', 'DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'purchase_canceled', 'DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'purchase_processed', 'DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'declined', 'DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'standin', 'DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'address', 'DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'creditanalysis', 'DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'person', 'DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'phone', 'DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'adjustment', 'DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'payment', 'DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'statement', 'DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'transaction', 'DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'bankslip', 'DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'account', 'DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'card', 'DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 242, 'tokenized_purchase', 'DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR', 'pier.hooks.DATAHUB_XPI_REGULAR_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'purchase_approved', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'purchase_cleared', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'limit', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'batch', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'purchase', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'purchase_canceled', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'purchase_processed', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'declined', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'standin', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'address', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'creditanalysis', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'person', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'phone', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'adjustment', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'payment', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'statement', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'transaction', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'bankslip', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'account', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'card', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 86, 'tokenized_purchase', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'purchase_approved', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'purchase_cleared', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'limit', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'batch', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'purchase', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'purchase_canceled', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'purchase_processed', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'declined', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'standin', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'address', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'creditanalysis', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'person', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'phone', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'adjustment', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'payment', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'statement', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'transaction', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'bankslip', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'account', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'card', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

INSERT INTO DELIVERY_QUEUES
(id, issuer_id, event_type, exchange_name, queue_name, queue_key, enabled, created_at, updated_at, deleted_at)
VALUES(gen_random_uuid(), 14, 'tokenized_purchase', 'DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED', 'pier.hooks.DATAHUB_THROTTLED_KEY', true, now(), now(), NULL);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back