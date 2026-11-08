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
