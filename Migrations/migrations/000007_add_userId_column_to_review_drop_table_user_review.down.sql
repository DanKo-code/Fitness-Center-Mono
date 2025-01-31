ALTER TABLE review
    DROP CONSTRAINT fk_user_id,
    DROP COLUMN user_id;

CREATE TABLE user_review (
    user_id UUID NOT NULL,
    review_id UUID NOT NULL,
    PRIMARY KEY (user_id, review_id),
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE,
    FOREIGN KEY (review_id) REFERENCES review(id) ON DELETE CASCADE
);
