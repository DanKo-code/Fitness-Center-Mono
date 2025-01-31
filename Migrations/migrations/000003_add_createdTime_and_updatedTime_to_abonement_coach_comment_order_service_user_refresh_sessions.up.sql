ALTER TABLE abonement
    ADD COLUMN created_time TIMESTAMPTZ,
    ADD COLUMN updated_time TIMESTAMPTZ;

ALTER TABLE coach
    ADD COLUMN created_time TIMESTAMPTZ,
    ADD COLUMN updated_time TIMESTAMPTZ;

ALTER TABLE comment
    DROP COLUMN create_date,
    ADD COLUMN created_time TIMESTAMPTZ,
    ADD COLUMN updated_time TIMESTAMPTZ;

ALTER TABLE "order"
    ADD COLUMN created_time TIMESTAMPTZ,
    ADD COLUMN updated_time TIMESTAMPTZ;

ALTER TABLE refresh_sessions
    ADD COLUMN created_time TIMESTAMPTZ,
    ADD COLUMN updated_time TIMESTAMPTZ;

ALTER TABLE service
    ADD COLUMN created_time TIMESTAMPTZ,
    ADD COLUMN updated_time TIMESTAMPTZ;

ALTER TABLE "user"
    DROP COLUMN created_time,
    ADD COLUMN created_time TIMESTAMPTZ,
    ADD COLUMN updated_time TIMESTAMPTZ;