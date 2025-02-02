CREATE TYPE training_status AS ENUM ('booked', 'active', 'passed');

CREATE TABLE training (
    id UUID PRIMARY KEY,
    time_from TIMESTAMPTZ NOT NULL,
    time_until TIMESTAMPTZ NOT NULL,
    status training_status NOT NULL,
    coach_id UUID NOT NULL,
    client_id UUID NOT NULL,
    created_time TIMESTAMPTZ NOT NULL,
    updated_time TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (coach_id)
        REFERENCES coach(id)
        ON DELETE CASCADE,
    FOREIGN KEY (client_id)
        REFERENCES "user"(id)
        ON DELETE CASCADE
)