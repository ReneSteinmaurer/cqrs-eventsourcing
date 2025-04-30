CREATE TABLE nutzer_history
(
    id              SERIAL PRIMARY KEY,
    nutzer_id       TEXT  NOT NULL,
    event_type      TEXT  NOT NULL,
    event_timestamp TIMESTAMP WITH TIME ZONE DEFAULT now(),
    payload         JSONB NOT NULL
);
