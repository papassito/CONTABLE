package domain

import "time"

// AuditAnchor gestiona el anclaje criptográfico a redes externas (Web3) para el Ledger inmutable de auditoría.
type AuditAnchor struct {
	ID              string    `json:"id"`
	ChainSequenceID int64     `json:"chain_sequence_id"`
	RootHash        string    `json:"root_hash"`
	TargetNetwork   string    `json:"target_network"`
	TxHash          string    `json:"tx_hash"`
	AnchoredAtUTC   time.Time `json:"anchored_at_utc"`
}
