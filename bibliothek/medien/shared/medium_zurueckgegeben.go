package shared

import "time"

type MediumZurueckgegebenEvent struct {
	MediumId string
	NutzerId string
	Date     time.Time
}

func NewMediumZurueckgegebenEvent(mediumId, nutzerId string, date time.Time) MediumZurueckgegebenEvent {
	return MediumZurueckgegebenEvent{
		MediumId: mediumId,
		NutzerId: nutzerId,
		Date:     date,
	}
}

var MediumZurueckgegebenEventType = "MediumZurueckgegebenEvent"
