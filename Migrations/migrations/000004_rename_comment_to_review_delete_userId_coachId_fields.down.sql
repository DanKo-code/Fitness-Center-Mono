ALTER TABLE review RENAME TO comment;

ALTER TABLE comment
    ADD COLUMN user_id uuid,
    ADD COLUMN coach_id uuid;