CREATE TABLE user_review (
    user_id UUID NOT NULL,
    review_id UUID NOT NULL,
    PRIMARY KEY (user_id, review_id),
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE,
    FOREIGN KEY (review_id) REFERENCES review(id) ON DELETE CASCADE
);

CREATE TABLE coach_review (
    coach_id UUID NOT NULL,
    review_id UUID NOT NULL,
    PRIMARY KEY (coach_id, review_id),
    FOREIGN KEY (coach_id) REFERENCES coach(id) ON DELETE CASCADE,
    FOREIGN KEY (review_id) REFERENCES review(id) ON DELETE CASCADE
);