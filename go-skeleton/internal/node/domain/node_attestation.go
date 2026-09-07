package domain

import (
	"errors"
	"time"
)

// NodeAttestation maneja la atestación física de llaves de hardware
type NodeAttestation struct {
	ID         string    `json:"id"`
	NodeID     string    `json:"node_id"`
	Challenge  string    `json:"challenge"`
	Signature  string    `json:"signature"`
	AttestedAt time.Time `json:"attested_at"`
	IsVerified bool      `json:"is_verified"`
}

// Validate verifica la validez estructural de la atestación
func (na *NodeAttestation) Validate() error {
	if na.ID == "" {
		return errors.New("el ID de atestación es requerido")
	}
	if na.NodeID == "" {
		return errors.New("el ID de nodo es requerido")
	}
	if na.Challenge == "" {
		return errors.New("el reto de atestación no puede estar vacío")
	}
	return nil
}
