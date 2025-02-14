-- liquibase formatted sql
-- changeset author:vlad

create table transactions(
    id bigserial primary key,
    sender_id int not null references users (id),
    recipient_id int not null references users (id),
    amount int not null,
    occurred timestamp not null,
    check (amount > 0),
    check (sender_id != recipient_id)
);

-- rollback DROP TABLE transactions;