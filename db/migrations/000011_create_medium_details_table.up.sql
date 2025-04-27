CREATE TABLE medium_details
(
    medium_id              UUID PRIMARY KEY,
    isbn                   TEXT,
    titel                  TEXT,
    genre                  TEXT,
    typ                    TEXT,
    standort               TEXT,
    signatur               TEXT,
    exemplar_code          TEXT,
    aktuell_verliehen      BOOLEAN DEFAULT FALSE,
    verliehen_an           TEXT,
    verliehen_an_nutzer_id TEXT,
    verliehen_von          TIMESTAMP,
    faellig_bis            TIMESTAMP
);