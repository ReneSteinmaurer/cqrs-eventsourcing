package shared

import "time"

type ProjectionStatus struct {
	ProjectionName         string    `db:"projection_name"`
	LastProcessedEventId   string    `db:"last_processed_event_id"`
	LastProcessedTimestamp time.Time `db:"last_processed_timestamp"`
}
