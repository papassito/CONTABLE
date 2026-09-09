package domain

import "time"

// NodeAttestation maneja la atestación física de llaves de hardware
type NodeAttestation struct {
	ID         string    `json:"id"`
	NodeID     string    `json:"node_id"`
	Challenge  string    `json:"challenge"`
	Signature  string    `json:"signature"`
	AttestedAt time.Time `json:"attested_at"`
	IsVerified bool      `json:"is_verified"`
}
