-- name: GetCalls :many
select c.id, c.duration, c.created_at,
       a.id as subscriber_id, a.name, a.phone_number
from calls c
left join call_participants cp on cp.call_id = c.id
left join abonents a on a.id = cp.subscriber_id
order by c.id, a.id;

-- name: CreateCall :one
insert into calls (duration, created_at)
values ($1, $2)
returning *;

-- name: AddCallParticipants :exec
insert into call_participants (call_id, subscriber_id)
select $1, unnest($2::bigint[]);

-- name: GetCallParticipants :many
select a.*
from abonents a
join call_participants cp on cp.subscriber_id = a.id
where cp.call_id = $1
order by a.id;
