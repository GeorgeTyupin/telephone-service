-- +goose Up
create table abonents (
    id int PRIMARY KEY,
    name text,
    phone_number text
);

create table calls (
    id int PRIMARY KEY,
    first_abonent int references abonents(id),
    second_abonent int references abonents(id),
    duration interval,
    created_at timestamptz
);

-- +goose Down
drop table abonents;
drop table calls;