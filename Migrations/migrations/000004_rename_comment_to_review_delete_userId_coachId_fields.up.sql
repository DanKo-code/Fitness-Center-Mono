ALTER TABLE comment RENAME TO review;

ALTER TABLE review
    DROP COLUMN user_id,
    DROP COLUMN coach_id;