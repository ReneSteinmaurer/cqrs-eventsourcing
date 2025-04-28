package events

import "time"

const (
	MediumVerlorenDurchBenutzerEventType = "MediumVerlorenDurchBenutzerEvent"
)

type MediumVerlorenDurchBenutzerEvent struct {
	MediumId string
	NutzerId string
	Date     time.Time
}

func NewMediumVerlorenEvent(mediumId, nutzerId string, date time.Time) MediumVerlorenDurchBenutzerEvent {
	return MediumVerlorenDurchBenutzerEvent{
		MediumId: mediumId,
		NutzerId: nutzerId,
		Date:     date,
	}
}
