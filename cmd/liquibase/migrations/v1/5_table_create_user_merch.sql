-- liquibase formatted sql
-- changeset author:vlad

create table purchases(
    user_id int references users (id),
    merch_id int references merch (id),
    quantity int not null default 1,
    occurred timestamp not null,
    primary key (user_id, occurred, merch_id)
);

-- rollback DROP TABLE purchases;