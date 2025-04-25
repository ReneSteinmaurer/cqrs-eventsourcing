package events

import "time"

const (
	MediumVerliehenEventType = "MediumVerliehenEvent"
)

type MediumVerliehenEvent struct {
	MediumId string
	NutzerId string
	Von      time.Time
	Bis      time.Time
}

func NewMediumVerliehenEvent(mediumId, nutzerId string, von, bis time.Time) MediumVerliehenEvent {
	return MediumVerliehenEvent{
		MediumId: mediumId,
		NutzerId: nutzerId,
		Von:      von,
		Bis:      bis,
	}
}
