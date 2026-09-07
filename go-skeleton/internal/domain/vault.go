package domain

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"io"
)

// EncryptedEnvelope representa la estructura de almacenamiento seguro de un secreto fiscal.
type EncryptedEnvelope struct {
	EncryptedPayload   []byte
	WrappedTEKNode     []byte
	WrappedTEKRecovery []byte
	PayloadNonce       []byte
}

// VaultService gestiona el cifrado híbrido y recuperación de secretos.
type VaultService struct {
	nodePublicKey     *rsa.PublicKey
	nodePrivateKey    *rsa.PrivateKey // Solo disponible en el nodo autorizado
	recoveryPublicKey *rsa.PublicKey  // Llave pública maestra (Cold Storage)
}

// NewVaultService inicializa el servicio de criptografía con las llaves RSA pertinentes.
func NewVaultService(nodePub *rsa.PublicKey, nodePriv *rsa.PrivateKey, recPub *rsa.PublicKey) *VaultService {
	return &VaultService{
		nodePublicKey:     nodePub,
		nodePrivateKey:    nodePriv,
		recoveryPublicKey: recPub,
	}
}

// SealSecret cifra un secreto en texto plano y genera los sobres digitales para el nodo y recuperación.
func (v *VaultService) SealSecret(plainSecret []byte) (*EncryptedEnvelope, error) {
	// 1. Generar TEK (Tenant Encryption Key) AES-256 (32 bytes)
	tek := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, tek); err != nil {
		return nil, errors.New("fallo al generar TEK")
	}

	// 2. Cifrar el secreto usando AES-GCM
	block, err := aes.NewCipher(tek)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.New("fallo al generar nonce")
	}

	encryptedPayload := gcm.Seal(nil, nonce, plainSecret, nil)

	// 3. Empaquetar (Wrap) la TEK para el Nodo
	wrappedTEKNode, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, v.nodePublicKey, tek, nil)
	if err != nil {
		return nil, errors.New("fallo al empaquetar llave de nodo")
	}

	// 4. Empaquetar (Wrap) la TEK para la Llave Maestra
	wrappedTEKRecovery, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, v.recoveryPublicKey, tek, nil)
	if err != nil {
		return nil, errors.New("fallo al empaquetar llave maestra")
	}

	return &EncryptedEnvelope{
		EncryptedPayload:   encryptedPayload,
		WrappedTEKNode:     wrappedTEKNode,
		WrappedTEKRecovery: wrappedTEKRecovery,
		PayloadNonce:       nonce,
	}, nil
}

// UnsealSecretWithNodeKey intenta descifrar el secreto usando la llave privada local del nodo.
func (v *VaultService) UnsealSecretWithNodeKey(env *EncryptedEnvelope) ([]byte, error) {
	if v.nodePrivateKey == nil {
		return nil, errors.New("llave privada de nodo no disponible")
	}

	// 1. Desempaquetar la TEK
	tek, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, v.nodePrivateKey, env.WrappedTEKNode, nil)
	if err != nil {
		return nil, errors.New("desempaquetado RSA denegado o llave incorrecta")
	}

	// 2. Descifrar el secreto AES-GCM
	block, err := aes.NewCipher(tek)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	decryptedSecret, err := gcm.Open(nil, env.PayloadNonce, env.EncryptedPayload, nil)
	if err != nil {
		return nil, errors.New("fallo de integridad al abrir AES-GCM")
	}

	return decryptedSecret, nil
}
