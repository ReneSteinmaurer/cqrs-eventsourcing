export type Ausleihstatus = 'AKTIV' | 'ÜBERFÄLLIG' | 'BALD_FÄLLIG';
export type NutzerStatus = 'AKTIV' | 'GESPERRT';

export interface AktiveAusleihe {
  mediumId: string;
  titel: string;
  ausgeliehenAm: Date | null;
  faelligAm: Date | null;
  status: Ausleihstatus;
}

export interface VerlorenesMedium {
  mediumId: string;
  titel: string;
  ausgeliehenAm: Date | null;
  faelligAm: Date | null;
}

export interface Notiz {
  text: string;
  erstelltAm: Date | null;
  erstelltVon: Date | null;
}

export interface NutzerDetails {
  nutzerId: string;
  vorname: string;
  nachname: string;
  status: NutzerStatus;
  registriertAm: Date | null;

  aktiveAusleihen: AktiveAusleihe[];
  letzteNotizen: Notiz[];
  verloreneMedien: VerlorenesMedium[];

  sperrgrund: string | null;
}
