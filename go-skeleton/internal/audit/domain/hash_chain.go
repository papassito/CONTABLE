package domain

import "time"

// AuditHashChainNode representa un bloque verificado en el Ledger inmutable del subsistema de auditoría.
type AuditHashChainNode struct {
	SequenceID           int64     `json:"sequence_id"`
	EventID              string    `json:"event_id"`
	TenantID             string    `json:"tenant_id"`
	NodeID               string    `json:"node_id"`
	ActorID              string    `json:"actor_id"`
	EventType            string    `json:"event_type"`
	CanonicalPayloadHash string    `json:"canonical_payload_hash"`
	PreviousHash         string    `json:"previous_hash"`
	ChainHash            string    `json:"chain_hash"`
	CreatedAtUTC         time.Time `json:"created_at_utc"`
}
