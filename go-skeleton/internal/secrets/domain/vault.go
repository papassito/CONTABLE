package domain

import "context"

// SecretsVault define el contrato para el cifrado y manejo de secretos.
// Se ubica en el Bounded Context de 'secrets' para cohesión.
type SecretsVault interface {
	// Encrypt toma un texto plano, lo cifra usando una TEK (Tenant Encryption Key),
	// y envuelve la TEK usando las llaves del nodo y de recuperación.
	Encrypt(ctx context.Context, tenantID, secretType string, plaintext []byte) (*SecretEnvelope, error)
	// Decrypt usa la llave privada del nodo para desenvolver la TEK y luego descifra el payload.
	Decrypt(ctx context.Context, envelope *SecretEnvelope) ([]byte, error)
}
