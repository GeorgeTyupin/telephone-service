-- +goose Up
create table call_participants (
    call_id bigint not null references calls(id) on delete cascade,
    subscriber_id bigint not null references abonents(id) on delete restrict,
    primary key (call_id, subscriber_id)
);

create index call_participants_subscriber_id_idx
    on call_participants (subscriber_id);

alter table calls
    drop column first_abonent,
    drop column second_abonent;

-- +goose Down
alter table calls
    add column first_abonent integer references abonents(id),
    add column second_abonent integer references abonents(id);

drop table call_participants;
