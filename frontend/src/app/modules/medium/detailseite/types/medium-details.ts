import { HistoryEvent } from "../../../../shared/types/history-events";

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
