-- liquibase formatted sql
-- changeset author:vlad

create table users(
    id serial primary key,
    name varchar(32) not null unique,
    passhash varchar(64) not null,
    coins int not null default 0,
    check (coins >= 0)
);

-- rollback DROP TABLE users;