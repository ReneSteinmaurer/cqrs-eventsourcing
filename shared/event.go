package shared

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	Id            string
	Type          string
	Timestamp     time.Time
	Payload       []byte
	Version       int
	AggregateType string
	AggregateKey  string
}

func NewEvent(aggregateType, aggregateKey, eventType string, version int, payloadJSON []byte) Event {
	return Event{
		Id:            uuid.NewString(),
		Type:          eventType,
		Timestamp:     time.Now().UTC(),
		Payload:       payloadJSON,
		AggregateKey:  aggregateKey,
		AggregateType: aggregateType,
		Version:       version,
	}
}
