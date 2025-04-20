create table events
(
    id             UUID primary key,
    type           text      not null,
    timestamp      timestamp not null,
    payload        JSONB     not null,
    aggregate_key  text      not null,
    aggregate_type text      not null,
    version        integer   not null,

    unique (aggregate_type, aggregate_key, version)
);

CREATE INDEX idx_events_aggregate_id ON events (aggregate_key);