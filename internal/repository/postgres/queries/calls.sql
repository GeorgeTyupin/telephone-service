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
where (c.id = @call_id::bigint or @call_id::bigint = 0)
  and (c.status = @status::text or @status::text = '')
  and (exists (
      select 1 from call_participants participant
      where participant.call_id = c.id
        and participant.subscriber_id = @subscriber_id::bigint
  ) or @subscriber_id::bigint = 0)
order by c.id, a.id, p.id;

-- name: StartCall :one
insert into calls (duration, created_at, status, started_at)
values (interval '0', statement_timestamp(), 'active', statement_timestamp())
returning *;

-- name: AddCallParticipants :exec
insert into call_participants (call_id, subscriber_id)
select @call_id, unnest(@subscriber_ids::bigint[])
on conflict (call_id, subscriber_id) do nothing;

-- name: StartCallParticipations :exec
insert into call_participations (call_id, subscriber_id, joined_at)
select @call_id, unnest(@subscriber_ids::bigint[]), @joined_at;

-- name: UpdateCall :one
update calls
set status = @status,
    duration = @duration,
    finished_at = @finished_at,
    version = version + 1
where id = @call_id and version = @version
returning id;

-- name: CloseCallParticipations :exec
update call_participations
set left_at = @left_at
where call_id = @call_id
  and id = any(@participation_ids::bigint[])
  and left_at is null;
