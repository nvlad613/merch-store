-- liquibase formatted sql
-- changeset author:vlad

CREATE INDEX t_sender_idx ON transactions USING HASH (sender_id);
CREATE INDEX t_recipient_idx ON transactions USING HASH (recipient_id);

-- rollback DROP INDEX t_sender_idx; DROP INDEX t_sender_idx t_recipient_idx;