-- +goose Up
alter table calls
    add column status text not null default 'finished',
    add column started_at timestamptz,
    add column finished_at timestamptz,
    add column version bigint not null default 1,
    add constraint calls_status_check check (status in ('active', 'finished')),
    add constraint calls_runtime_check check (
        (status = 'active' and started_at is not null and finished_at is null)
        or (status = 'finished' and (
            (started_at is null and finished_at is null)
            or (started_at is not null and finished_at is not null
                and finished_at >= started_at)
        ))
    );

create table call_participations (
    id bigint generated always as identity primary key,
    call_id bigint not null,
    subscriber_id bigint not null,
    joined_at timestamptz not null,
    left_at timestamptz,
    foreign key (call_id, subscriber_id)
        references call_participants (call_id, subscriber_id) on delete cascade,
    check (left_at is null or left_at >= joined_at)
);

create unique index call_participations_connected_key
    on call_participations (call_id, subscriber_id) where left_at is null;

create index call_participations_history_idx
    on call_participations (call_id, subscriber_id, id);

-- +goose Down
drop table call_participations;

alter table calls
    drop column version,
    drop column finished_at,
    drop column started_at,
    drop column status;
