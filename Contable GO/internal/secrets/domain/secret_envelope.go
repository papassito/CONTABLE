package domain

import "time"

// SecretEnvelope encapsula credenciales con cifrado multicapa AES-256-GCM y RSA.
// Este archivo centraliza la definición del modelo de dominio para la bóveda de secretos.
type SecretEnvelope struct {
	ID                 string    `json:"id"`
	TenantID           string    `json:"tenant_id"`
	SecretType         string    `json:"secret_type"`
	EncryptedPayload   []byte    `json:"encrypted_payload"`
	WrappedTekNode     []byte    `json:"wrapped_tek_node"`
	WrappedTekRecovery []byte    `json:"wrapped_tek_recovery"`
	PayloadNonce       []byte    `json:"payload_nonce"`
	KeyVersion         string    `json:"key_version"`
	CreatedAtUTC       time.Time `json:"created_at_utc"`
	UpdatedAtUTC       time.Time `json:"updated_at_utc"`
}
