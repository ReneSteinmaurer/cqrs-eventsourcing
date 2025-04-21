package shared

type MediumKatalogisiertEvent struct {
	MediumId     string
	Signature    string
	Standort     string
	ExemplarCode string
}

func NewMediumKatalogisiertEvent(mediumId, signature, standort, exemplarCode string) MediumKatalogisiertEvent {
	return MediumKatalogisiertEvent{
		MediumId:     mediumId,
		Signature:    signature,
		Standort:     standort,
		ExemplarCode: exemplarCode,
	}
}

var MediumKatalogisiertEventType = "MediumKatalogisiertEvent"
