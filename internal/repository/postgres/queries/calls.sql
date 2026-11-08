-- name: GetCalls :many
select c.id, c.duration, c.created_at, c.status,
       c.started_at, c.finished_at, c.version,
       statement_timestamp()::timestamptz as observed_at,
       (case when c.status = 'active'
           then greatest(statement_timestamp() - c.started_at, interval '0')
           else c.duration end)::interval as elapsed,
       a.id as subscriber_id, a.name, a.phone_number,
       p.id as participation_id, p.joined_at, p.left_at
from calls c
left join call_participants cp on cp.call_id = c.id
left join abonents a on a.id = cp.subscriber_id
left join call_participations p
    on p.call_id = cp.call_id and p.subscriber_id = cp.subscriber_id
where (c.id = $1::bigint or $1::bigint = 0)
  and (c.status = $2::text or $2::text = '')
  and (exists (
      select 1 from call_participants participant
      where participant.call_id = c.id
        and participant.subscriber_id = $3::bigint
  ) or $3::bigint = 0)
order by c.id, a.id, p.id;

-- name: StartCall :one
insert into calls (duration, created_at, status, started_at)
values (interval '0', statement_timestamp(), 'active', statement_timestamp())
returning *;

-- name: AddCallParticipants :exec
insert into call_participants (call_id, subscriber_id)
select $1, unnest($2::bigint[])
on conflict (call_id, subscriber_id) do nothing;

-- name: StartCallParticipations :exec
insert into call_participations (call_id, subscriber_id, joined_at)
select $1, unnest($2::bigint[]), $3;

-- name: UpdateCall :one
update calls
set status = $1,
    duration = $2,
    finished_at = $3,
    version = version + 1
where id = $4 and version = $5
returning id;

-- name: CloseCallParticipations :exec
update call_participations
set left_at = $1
where call_id = $2
  and id = any($3::bigint[])
  and left_at is null;
