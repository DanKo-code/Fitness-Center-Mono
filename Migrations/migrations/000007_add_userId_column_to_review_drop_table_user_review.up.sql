ALTER TABLE review
    ADD COLUMN user_id uuid,
    ADD CONSTRAINT fk_user_id
    FOREIGN KEY (user_id) REFERENCES "user"(id);

DROP TABLE user_review;
