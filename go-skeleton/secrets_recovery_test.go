package fcos_kernel_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io"
	"testing"
)

func TestSecretsVaultKeyEnvelopeAndMasterRecovery(t *testing.T) {
	// 1. Generar pares RSA para el Nodo A y la Clave de Recuperación Maestra
	nodeKeypairA, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Error generando RSA Node A: %v", err)
	}

	masterRecoveryKeypair, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Error generando RSA Recovery: %v", err)
	}

	// 2. Generar TEK (Tenant Encryption Key) simétrica de 32 bytes (AES-256)
	tek := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, tek); err != nil {
		t.Fatalf("Error generando TEK: %v", err)
	}

	// 3. Cifrar Secreto Fiscal (.key de e.firma) con TEK usando AES-GCM
	fiscalSecret := []byte("LLAVE_PRIVADA_EFIRMA_SECRET_PAYLOAD_CONTRIBUYENTE")
	block, err := aes.NewCipher(tek)
	if err != nil {
		t.Fatalf("Error cipher AES: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatalf("Error GCM: %v", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	encryptedPayload := gcm.Seal(nil, nonce, fiscalSecret, nil)

	// 4. Crear Key Envelopes para Nodo A y para Máster de Recuperación
	wrappedTekNode, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &nodeKeypairA.PublicKey, tek, nil)
	if err != nil {
		t.Fatalf("Error Key Wrapping Node: %v", err)
	}

	wrappedTekRecovery, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &masterRecoveryKeypair.PublicKey, tek, nil)
	if err != nil {
		t.Fatalf("Error Key Wrapping Recovery: %v", err)
	}

	// 5. Caso Normal: Descifrar TEK usando la clave del Nodo A
	unwrappedTekNode, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, nodeKeypairA, wrappedTekNode, nil)
	if err != nil || !bytes.Equal(unwrappedTekNode, tek) {
		t.Fatalf("FALLO DE DESCIFRADO NODO: La TEK recuperada no coincide")
	}

	// 6. Caso de Desastre / Cambio de Hardware: Reemplazo de Placa Madre (Nuevo Keypair Node B)
	nodeKeypairB, _ := rsa.GenerateKey(rand.Reader, 2048)
	_, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, nodeKeypairB, wrappedTekNode, nil)
	if err == nil {
		t.Fatalf("FALLO DE SEGURIDAD: Un nodo no autorizado pudo descifrar el Key Envelope")
	}

	// 7. Fallback de Recuperación: Descifrar TEK usando la Clave Maestra de Recuperación
	unwrappedTekRecovery, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, masterRecoveryKeypair, wrappedTekRecovery, nil)
	if err != nil || !bytes.Equal(unwrappedTekRecovery, tek) {
		t.Fatalf("FALLO DE RECUPERACIÓN MAESTRA: No se pudo restaurar la TEK vía Master Key Envelope")
	}

	// 8. Verificar descifrado final del secreto fiscal con la TEK recuperada
	blockRec, _ := aes.NewCipher(unwrappedTekRecovery)
	gcmRec, _ := cipher.NewGCM(blockRec)
	decryptedSecret, err := gcmRec.Open(nil, nonce, encryptedPayload, nil)
	if err != nil || !bytes.Equal(decryptedSecret, fiscalSecret) {
		t.Fatalf("CORRUPCIÓN DE SECRETO FISCAL: El secreto descifrado no coincide con el original")
	}
}
