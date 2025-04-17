package shared

import "time"

type Event struct {
	Id        string
	Type      string
	Timestamp time.Time
	Payload   []byte
}
