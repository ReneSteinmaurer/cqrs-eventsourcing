CREATE TABLE medium_historie
(
    id              SERIAL PRIMARY KEY,
    medium_id       UUID  NOT NULL,
    event_type      TEXT  NOT NULL,
    event_timestamp TIMESTAMP WITH TIME ZONE DEFAULT now(),
    payload         JSONB NOT NULL
);
