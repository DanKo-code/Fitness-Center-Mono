ALTER TABLE "user"
    DROP COLUMN first_name,
    DROP COLUMN last_name,
    ADD COLUMN name VARCHAR(255) NOT NULL;