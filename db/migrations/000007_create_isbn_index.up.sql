create table isbn_index(
    isbn text primary key ,
    medium_id text not null,
    created_at timestamp not null default now()
)