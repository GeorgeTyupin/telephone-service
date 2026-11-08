-- name: GetSubscribers :many
select id, name, phone_number
from abonents
order by id;

-- name: CreateSubscribers :exec
insert into abonents (name, phone_number)
select requested.name, requested.phone_number
from (
    select unnest(@names::text[]) as name,
           unnest(@phone_numbers::text[]) as phone_number
) requested
order by ltrim(regexp_replace(phone_number, '[^0-9]', '', 'g'), '0')
on conflict ((ltrim(regexp_replace(phone_number, '[^0-9]', '', 'g'), '0'))) do nothing;

-- name: GetSubscribersByPhoneNumbers :many
select a.id, a.name, a.phone_number
from abonents a
where ltrim(regexp_replace(a.phone_number, '[^0-9]', '', 'g'), '0') in (
    select ltrim(regexp_replace(requested.phone_number, '[^0-9]', '', 'g'), '0')
    from unnest(@phone_numbers::text[]) as requested(phone_number)
)
order by a.id
for key share of a;

-- name: GetSubscriber :one
select id, name, phone_number
from abonents
where id = @id;

-- name: GetSubscriberStats :one
select count(c.id)::bigint as calls_count,
       coalesce(sum(c.duration), interval '0')::interval as total_duration,
       coalesce(avg(c.duration), interval '0')::interval as average_duration
from abonents a
left join call_participants cp on cp.subscriber_id = a.id
left join calls c on c.id = cp.call_id and c.status = 'finished'
where a.id = @id
group by a.id;
