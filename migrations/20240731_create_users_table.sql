-- +goose Up
create table users (
    id serial primary key,
    name text not null UNIQUE,
    email text not null UNIQUE,
    password text not null,
    role int not null,
    created_at timestamp not null default now(),
    updated_at timestamp
);

-- +goose Down
drop table users;