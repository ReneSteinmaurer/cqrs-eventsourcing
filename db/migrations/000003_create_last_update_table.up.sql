create table projection_status (
    projection_name text primary key,
    last_processed_event_id uuid,
    last_processed_timestamp timestamp
)