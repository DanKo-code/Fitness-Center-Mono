CREATE TYPE training_status AS ENUM ('booked', 'active', 'passed');

CREATE TABLE training (
    id UUID PRIMARY KEY,
    "date" DATE NOT NULL,
    hour_from INT NOT NULL,
    hour_until INT NOT NULL,
    status training_status NOT NULL DEFAULT 'booked',
    coach_id UUID NOT NULL,
    client_id UUID NOT NULL,
    created_time TIMESTAMPTZ,
    updated_time time,
    FOREIGN KEY (coach_id)
        REFERENCES coach(id)
        ON DELETE CASCADE
)