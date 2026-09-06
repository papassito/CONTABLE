package domain

import "context"

// SecretsVault define el comportamiento esperado de la bóveda de llaves del nodo
type SecretsVault interface {
	Encrypt(ctx context.Context, tenantID, secretType string, plaintext []byte) (*SecretEnvelope, error)
	Decrypt(ctx context.Context, envelope *SecretEnvelope) ([]byte, error)
	LoadNodeKeys(ctx context.Context) error
}
