package events

import "time"

const (
	MediumVerlorenDurchBenutzerEventType = "MediumVerlorenDurchBenutzerEvent"
	MediumBestandsverlustEventType       = "MediumBestandsverlustEvent"
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

type MediumBestandsverlustEvent struct {
	MediumId string
	Date     time.Time
}

func NewMediumBestandsverlustEvent(mediumId string, date time.Time) MediumBestandsverlustEvent {
	return MediumBestandsverlustEvent{
		MediumId: mediumId,
		Date:     date,
	}
}
