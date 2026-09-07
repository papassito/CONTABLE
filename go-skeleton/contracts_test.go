package main_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/gowebpki/jcs"
)

func TestCanonicalJSONRFC8785Determinism(t *testing.T) {
	// Payload 1: Propiedades en orden A
	jsonA := []byte(`{
		"tenant_id": "a0000000-0000-0000-0000-000000000001",
		"base_amount_cents": 428518,
		"rule_code": "IMSS_RCV_V1",
		"expediente_id": "b0000000-0000-0000-0000-000000000002"
	}`)

	// Payload 2: Mismos datos, diferente orden e identación
	jsonB := []byte(`{
		"expediente_id": "b0000000-0000-0000-0000-000000000002",
		"rule_code": "IMSS_RCV_V1",
		"tenant_id": "a0000000-0000-0000-0000-000000000001",
		"base_amount_cents": 428518
	}`)

	// Canonicalización bajo RFC 8785
	canonicalA, err := jcs.Transform(jsonA)
	if err != nil {
		t.Fatalf("Error al canonicalizar jsonA: %v", err)
	}

	canonicalB, err := jcs.Transform(jsonB)
	if err != nil {
		t.Fatalf("Error al canonicalizar jsonB: %v", err)
	}

	// Cómputo de Hashes SHA-256
	hashA := sha256.Sum256(canonicalA)
	hashB := sha256.Sum256(canonicalB)

	hashAStr := hex.EncodeToString(hashA[:])
	hashBStr := hex.EncodeToString(hashB[:])

	if hashAStr != hashBStr {
		t.Fatalf("FALLO DE CANONICALIZACIÓN: Hashes discrepantes para payloads idénticos.\nHash A: %s\nHash B: %s", hashAStr, hashBStr)
	}
}
