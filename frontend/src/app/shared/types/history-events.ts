export const eventTypeKeys = [
  'NutzerRegistriertEvent',
  'MediumErworbenEvent',
  'MediumKatalogisiertEvent',
  'MediumVerliehenEvent',
  'MediumZurueckgegebenEvent',
  'MediumVerlorenDurchBenutzerEvent',
  'MediumBestandsverlustEvent',
  'MediumWiederaufgefundenEvent',
  'MediumWiederaufgefundenDurchNutzerEvent',
] as const

export const eventTypeToIcon: Record<keyof EventTypeToConcrete, string> = {
  NutzerRegistriertEvent: 'person_add',
  MediumErworbenEvent: 'label_important',
  MediumKatalogisiertEvent: 'library_add',
  MediumVerliehenEvent: 'logout',
  MediumZurueckgegebenEvent: 'keyboard_return',
  MediumVerlorenDurchBenutzerEvent: 'report',
  MediumWiederaufgefundenEvent: 'undo',
  MediumBestandsverlustEvent: 'cancel_presentation',
  MediumWiederaufgefundenDurchNutzerEvent: 'keyboard_return',
};

export type EventTypeToConcrete = {
  NutzerRegistriertEvent: NutzerRegistriertEvent;
  MediumErworbenEvent: MediumErworbenEvent;
  MediumKatalogisiertEvent: MediumKatalogisiertEvent;
  MediumVerliehenEvent: MediumVerliehenEvent;
  MediumZurueckgegebenEvent: MediumZurueckgegebenEvent;
  MediumVerlorenDurchBenutzerEvent: MediumVerlorenDurchBenutzerEvent;
  MediumWiederaufgefundenEvent: MediumWiederaufgefundenEvent;
  MediumBestandsverlustEvent: MediumBestandsverlustEvent;
  MediumWiederaufgefundenDurchNutzerEvent: MediumWiederaufgefundenDurchNutzerEvent;
};

export type EventType = typeof eventTypeKeys[number]

interface BaseHistoryEvent<TPayload> {
  eventType: EventType;
  timestamp: string;
  payload: TPayload;
}

interface NutzerRegistriertPayload {
  Vorname: string
  Nachname: string
  Email: string
  NutzerId: string
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

export type NutzerRegistriertEvent = BaseHistoryEvent<NutzerRegistriertPayload> & {
  eventType: 'NutzerRegistriertEvent';
};

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
  NutzerRegistriertEvent |
  MediumErworbenEvent |
  MediumKatalogisiertEvent |
  MediumVerliehenEvent |
  MediumZurueckgegebenEvent |
  MediumVerlorenDurchBenutzerEvent |
  MediumWiederaufgefundenEvent |
  MediumBestandsverlustEvent |
  MediumWiederaufgefundenDurchNutzerEvent;

export function isEventType<K extends keyof EventTypeToConcrete>(
  event: HistoryEvent,
  type: K
): event is EventTypeToConcrete[K] {
  return event.eventType === type;
}

