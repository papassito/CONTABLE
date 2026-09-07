package domain

import "context"
import "time"

// TransactionalOutboxEvent asegura la entrega de eventos bajo consistencia eventual en el módulo de mensajería.
type TransactionalOutboxEvent struct {
	ID            string    `json:"id"`
	AggregateID   string    `json:"aggregate_id"`
	AggregateType string    `json:"aggregate_type"`
	EventType     string    `json:"event_type"`
	PayloadJSON   string    `json:"payload_json"`
	IsPublished   bool      `json:"is_published"`
	CreatedAtUTC  time.Time `json:"created_at_utc"`
}

// OutboxRepository define el contrato de persistencia para el patrón Transactional Outbox.
type OutboxRepository interface {
	SaveEvent(ctx context.Context, event *TransactionalOutboxEvent) error
	GetUnpublishedEvents(ctx context.Context, limit int) ([]*TransactionalOutboxEvent, error)
	MarkAsPublished(ctx context.Context, eventID string) error
}
