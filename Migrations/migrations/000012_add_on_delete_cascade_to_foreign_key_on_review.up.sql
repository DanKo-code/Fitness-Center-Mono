ALTER TABLE "review" DROP CONSTRAINT IF EXISTS fk_user_id;

ALTER TABLE "review"
    ADD CONSTRAINT fk_user_id
        FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE;
