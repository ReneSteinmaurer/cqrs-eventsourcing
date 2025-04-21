create table medium_bestand
(
    medium_id     TEXT PRIMARY KEY,
    isbn          TEXT NOT NULL,
    medium_type   TEXT NOT NULL,
    name          TEXT NOT NULL,
    genre         TEXT,
    signature     TEXT,
    standort      TEXT,
    exemplar_code TEXT
)