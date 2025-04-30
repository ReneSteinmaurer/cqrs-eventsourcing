const MediumStatusKeys = [
  'ERWORBEN',
  'KATALOGISIERT',
  'VERLIEHEN',
  'VERLOREN'
] as const

export type MediumStatus = typeof MediumStatusKeys[number]

export interface MediumDetail {
  mediumId: string;
  isbn: string;
  titel: string;
  genre: string;
  typ: string;
  standort: string;
  signatur: string;
  exemplarCode: string;
  aktuellVerliehen: boolean;
  verliehenAn: string | null;
  verliehenNutzerId: string | null;
  verliehenVon: Date | null;
  verlorenVonNutzerId: string | null;
  verlorenAm: Date | null;
  verlorenNutzerName: string | null
  faelligBis: Date | null;
  erworbenAm: Date;
  katalogisiertAm: Date;
  status: MediumStatus;
  historie: HistoryEvent[];
}

interface BaseHistoryEvent<TPayload> {
  eventType: string;
  timestamp: string;
  payload: TPayload;
}

interface MediumErworbenPayload {
  ISBN: string;
  Name: string;
  Genre: string;
  MediumId: string;
  MediumType: string;
}

interface MediumKatalogisiertPayload {
  MediumId: string;
  Standort: string;
  Signature: string;
  ExemplarCode: string;
}

interface MediumVerliehenPayload {
  MediumId: string;
  Von?: string;
  Bis?: string;
  Date?: Date;
  NutzerId: string;
}

interface MediumZurueckgegebenPayload {
  MediumId: string;
  NutzerId: string;
  Date?: Date;
}

interface MediumVerlorenDurchNutzerPayload {
  MediumId: string;
  NutzerId: string;
  Date: Date;
}

interface MediumBestandsverlustPayload {
  MediumId: string;
  Date: Date;
}

interface MediumWiederaufgefundenPayload {
  MediumId: string;
  Date: Date;
}

interface MediumWiederaufgefundenPayload {
  MediumId: string;
  Date: Date;
}

interface MediumWiederaufgefundenDurchNutzerPayload {
  MediumId: string;
  NutzerId: string;
  Date: Date;
}

export type MediumErworbenEvent = BaseHistoryEvent<MediumErworbenPayload> & {
  eventType: 'MediumErworbenEvent';
};

export type MediumKatalogisiertEvent = BaseHistoryEvent<MediumKatalogisiertPayload> & {
  eventType: 'MediumKatalogisiertEvent';
};

export type MediumVerliehenEvent = BaseHistoryEvent<MediumVerliehenPayload> & {
  eventType: 'MediumVerliehenEvent';
};

export type MediumZurueckgegebenEvent = BaseHistoryEvent<MediumZurueckgegebenPayload> & {
  eventType: 'MediumZurueckgegebenEvent';
};

export type MediumVerlorenDurchBenutzerEvent = BaseHistoryEvent<MediumVerlorenDurchNutzerPayload> & {
  eventType: 'MediumVerlorenDurchBenutzerEvent';
};

export type MediumBestandsverlustEvent = BaseHistoryEvent<MediumBestandsverlustPayload> & {
  eventType: 'MediumBestandsverlustEvent';
};

export type MediumWiederaufgefundenEvent = BaseHistoryEvent<MediumWiederaufgefundenPayload> & {
  eventType: 'MediumWiederaufgefundenEvent';
};

export type MediumWiederaufgefundenDurchNutzerEvent = BaseHistoryEvent<MediumWiederaufgefundenDurchNutzerPayload> & {
  eventType: 'MediumWiederaufgefundenDurchNutzerEvent';
};

export type HistoryEvent =
  MediumErworbenEvent |
  MediumKatalogisiertEvent |
  MediumVerliehenEvent |
  MediumZurueckgegebenEvent |
  MediumVerlorenDurchBenutzerEvent |
  MediumWiederaufgefundenEvent |
  MediumBestandsverlustEvent |
  MediumWiederaufgefundenDurchNutzerEvent;
