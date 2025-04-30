package detailseite

import "time"

type Ausleihstatus string
type NutzerStatus string

const (
	AusleihstatusAktiv        Ausleihstatus = "AKTIV"
	AusleihstatusUeberfaellig Ausleihstatus = "ÜBERFÄLLIG"
	AusleihstatusBaldFaellig  Ausleihstatus = "BALD_FÄLLIG"

	NutzerstatusAktiv    NutzerStatus = "AKTIV"
	NutzerstatusGesperrt NutzerStatus = "GESPERRT"
)

type AktiveAusleihe struct {
	MediumId      string        `json:"mediumId"`
	Titel         string        `json:"titel"`
	AusgeliehenAm *time.Time    `json:"ausgeliehenAm"`
	FaelligAm     *time.Time    `json:"faelligAm"`
	Status        Ausleihstatus `json:"status"`
}

type VerlorenesMedium struct {
	MediumId      string     `json:"mediumId"`
	Titel         string     `json:"titel"`
	AusgeliehenAm *time.Time `json:"ausgeliehenAm"`
	FaelligAm     *time.Time `json:"faelligAm"`
}

type Notiz struct {
	Text        string     `json:"text"`
	ErstelltAm  *time.Time `json:"erstelltAm"`
	ErstelltVon *time.Time `json:"erstelltVon"`
}

type NutzerDetails struct {
	NutzerId      string       `json:"nutzerId"`
	Vorname       string       `json:"vorname"`
	Nachname      string       `json:"nachname"`
	Status        NutzerStatus `json:"status"`
	RegistriertAm *time.Time   `json:"registriertAm"`

	AktiveAusleihen []AktiveAusleihe   `json:"aktiveAusleihen"`
	LetzteNotizen   []Notiz            `json:"letzteNotizen"`
	VerloreneMedien []VerlorenesMedium `json:"verloreneMedien"`

	Sperrgrund *string `json:"sperrgrund"`
}
