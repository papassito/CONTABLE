package domain

import "context"

// SecretsVault define el contrato para el cifrado y manejo de secretos.
type SecretsVault interface {
	Encrypt(ctx context.Context, tenantID string, plaintext []byte) (*SecretEnvelope, error)
}
