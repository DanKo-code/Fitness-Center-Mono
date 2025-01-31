ALTER TABLE abonement
    DROP COLUMN created_time,
    DROP COLUMN updated_time;

ALTER TABLE coach
    DROP COLUMN created_time,
    DROP COLUMN updated_time;

ALTER TABLE comment
    DROP COLUMN created_time,
    DROP COLUMN updated_time,
    ADD COLUMN create_date TIMESTAMP;

ALTER TABLE "order"
    DROP COLUMN created_time,
    DROP COLUMN updated_time;

ALTER TABLE refresh_sessions
    DROP COLUMN created_time,
    DROP COLUMN updated_time;

ALTER TABLE service
    DROP COLUMN created_time,
    DROP COLUMN updated_time;

ALTER TABLE "user"
    DROP COLUMN created_time,
    DROP COLUMN updated_time,
    ADD COLUMN created_time TIMESTAMP;
