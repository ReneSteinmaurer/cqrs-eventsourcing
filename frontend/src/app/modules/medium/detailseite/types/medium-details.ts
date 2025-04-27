const MediumStatusKeys = [
  'ERWORBEN',
  'KATALOGISIERT',
  'VERLIEHEN'
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

export type HistoryEvent = MediumErworbenEvent | MediumKatalogisiertEvent | MediumVerliehenEvent | MediumZurueckgegebenEvent;
