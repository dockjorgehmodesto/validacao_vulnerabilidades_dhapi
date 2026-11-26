-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

insert into issuer_events(name, connection_id, issuer_id, created_at, updated_at)
select a.event_type,
       a.id,
       a.issuer_id,
       now(),
       now()
from amqp_connections a
left join issuer_events ie on (a.issuer_id = ie.issuer_id and a.id = ie.connection_id and a.event_type = ie.name)
where ie.id is null;



-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
