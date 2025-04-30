create table nutzer_details
(
    nutzer_id         text      not null primary key,
    vorname           text      not null,
    nachname          text      not null,
    status            text      not null,
    registriert_am    timestamp not null,

    aktive_ausleihen  jsonb     not null default '[]',

    letzte_notizen    jsonb     not null default '[]',

    letzte_ereignisse jsonb     not null default '[]',

    sperrgrund        text
);