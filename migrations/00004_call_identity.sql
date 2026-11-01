-- +goose Up
alter table calls alter column id add generated always as identity;

-- +goose Down
alter table calls alter column id drop identity;
