package events

import "time"

const (
	MediumWiederaufgefundenEventType = "MediumWiederaufgefundenEvent"
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
