package domain

import "time"

// SecretEnvelope representa la estructura de un contenedor cifrado de secretos.
type SecretEnvelope struct {
	ID           string    `json:"id"`
	KeyID        string    `json:"key_id"`
	Ciphertext   []byte    `json:"ciphertext"`
	CreatedAtUTC time.Time `json:"created_at_utc"`
}
