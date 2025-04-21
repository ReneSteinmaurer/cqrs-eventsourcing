create table medium_verliehen
(
    medium_id     text primary key,
    verliehen_von timestamp not null,
    verliehen_bis timestamp not null,
    nutzer_id     text      not null
)