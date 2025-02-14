-- liquibase formatted sql
-- changeset author:vlad

create table merch(
    id serial primary key,
    name varchar(32) unique,
    price int not null check (price > 0)
);

-- rollback DROP TABLE merch;