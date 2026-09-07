package domain

import "time"

// TransactionalOutboxEvent asegura la entrega de eventos bajo consistencia eventual
type TransactionalOutboxEvent struct {
	ID            string    `json:"id"`
	AggregateID   string    `json:"aggregate_id"`
	AggregateType string    `json:"aggregate_type"`
	EventType     string    `json:"event_type"`
	PayloadJSON   string    `json:"payload_json"`
	IsPublished   bool      `json:"is_published"`
	CreatedAtUTC  time.Time `json:"created_at_utc"`
}
