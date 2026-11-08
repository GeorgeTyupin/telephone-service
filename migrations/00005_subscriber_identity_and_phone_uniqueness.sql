-- +goose Up
alter table abonents alter column id add generated always as identity;

select setval(
    pg_get_serial_sequence('abonents', 'id'),
    greatest(coalesce((select max(id) from abonents), 0), 0) + 1,
    false
);

create unique index abonents_phone_number_key
    on abonents ((ltrim(regexp_replace(phone_number, '[^0-9]', '', 'g'), '0')));

-- +goose Down
drop index abonents_phone_number_key;

alter table abonents alter column id drop identity;
