package shared

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	Id          string
	Type        string
	Timestamp   time.Time
	Payload     []byte
	AggregateId string
	Version     int
}

func NewEvent(aggregateId, eventType string, version int, payloadJSON []byte) Event {
	return Event{
		Id:          uuid.NewString(),
		Type:        eventType,
		Timestamp:   time.Now().UTC(),
		Payload:     payloadJSON,
		AggregateId: aggregateId,
		Version:     version,
	}
}
