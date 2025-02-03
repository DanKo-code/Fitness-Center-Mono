CREATE TYPE training_status AS ENUM ('booked', 'active', 'passed');

CREATE TABLE training (
    id UUID PRIMARY KEY,
    time_from timestamp NOT NULL,
    time_until timestamp NOT NULL,
    status training_status NOT NULL,
    coach_id UUID NOT NULL,
    client_id UUID NOT NULL,
    created_time timestamp NOT NULL,
    updated_time timestamp NOT NULL,
    FOREIGN KEY (coach_id)
        REFERENCES coach(id)
        ON DELETE CASCADE,
    FOREIGN KEY (client_id)
        REFERENCES "user"(id)
        ON DELETE CASCADE
)