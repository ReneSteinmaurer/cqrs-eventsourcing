package events

import "time"

const (
	MediumWiederaufgefundenEventType            = "MediumWiederaufgefundenEvent"
	MediumWiederaufgefundenDurchNutzerEventType = "MediumWiederaufgefundenDurchNutzerEvent"
)

type MediumWiederaufgefundenEvent struct {
	MediumId string
	Date     time.Time
}

func NewMediumWiederaufgefundenEvent(mediumId string, date time.Time) MediumWiederaufgefundenEvent {
	return MediumWiederaufgefundenEvent{
		MediumId: mediumId,
		Date:     date,
	}
}

type MediumWiederaufgefundenDurchNutzerEvent struct {
	MediumId string
	NutzerId string
	Date     time.Time
}

func NewMediumWiederaufgefundenDurchNutzerEvent(mediumId, nutzerId string, date time.Time) MediumWiederaufgefundenDurchNutzerEvent {
	return MediumWiederaufgefundenDurchNutzerEvent{
		MediumId: mediumId,
		Date:     date,
	}
}
