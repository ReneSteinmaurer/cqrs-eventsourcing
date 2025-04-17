create table events (
    id UUID primary key,
    type text not null,
    timestamp timestamp not null,
    payload JSONB not null
)