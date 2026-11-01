-- name: GetSubscribers :many
select id, name, phone_number
from abonents
order by id;
