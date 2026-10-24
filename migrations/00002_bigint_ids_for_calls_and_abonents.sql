-- +goose Up
alter table calls alter column id type bigint;
alter table abonents alter column id type bigint;

-- +goose Down
alter table calls alter column id type int;
alter table abonents alter column id type int;
