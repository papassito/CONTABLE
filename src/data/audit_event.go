package domain

import "time"

// AuditEvent representa una mutación canónica candidata a asentar en la cadena
type AuditEvent struct {
	EventID              string    `json:"event_id"`
	TenantID             string    `json:"tenant_id"`
	NodeID               string    `json:"node_id"`
	ActorID              string    `json:"actor_id"`
	EventType            string    `json:"event_type"`
	CanonicalPayloadHash string    `json:"canonical_payload_hash"`
	CreatedAtUTC         time.Time `json:"created_at_utc"`
}
