-- liquibase formatted sql
-- changeset author:vlad

create table user_merch(
    user_id int references users (id),
    merch_id int references merch (id),
    quantity int not null default 1,
    check (quantity > 0),
    primary key (user_id, merch_id)
);

-- rollback DROP TABLE user_merch;