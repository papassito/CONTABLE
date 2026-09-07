package domain

import "time"

// AuditEvent representa un registro de evento de auditoría dentro del subsistema de auditoría dedicado.
type AuditEvent struct {
	EventID              string    `json:"event_id"`
	TenantID             string    `json:"tenant_id"`
	NodeID               string    `json:"node_id"`
	ActorID              string    `json:"actor_id"`
	EventType            string    `json:"event_type"`
	CanonicalPayloadHash string    `json:"canonical_payload_hash"`
	CreatedAtUTC         time.Time `json:"created_at_utc"`
}
