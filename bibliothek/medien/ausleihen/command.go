package ausleihen

type VerleiheMediumCommand struct {
	MediumId string `json:"mediumId"`
	NutzerId string `json:"nutzerId"`
}
