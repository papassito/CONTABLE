package integration_test

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"testing"

	_ "modernc.org/sqlite"
)

func TestAuditChainTamperDetectionAndRejection(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Error al abrir base de datos: %v", err)
	}
	defer db.Close()

	schema := `
	CREATE TABLE audit_hash_chain (
		sequence_id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_id TEXT UNIQUE,
		tenant_id TEXT,
		node_id TEXT,
		actor_id TEXT,
		event_type TEXT,
		canonical_payload_hash TEXT,
		previous_hash TEXT,
		chain_hash TEXT,
		created_at_utc TEXT
	);

	CREATE TRIGGER trg_prevent_audit_update BEFORE UPDATE ON audit_hash_chain
	BEGIN
		SELECT RAISE(ABORT, 'VIOLACIÓN CONSTITUCIONAL: Inmutabilidad violada.');
	END;
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Error al crear esquema: %v", err)
	}

	// Insertar Evento Génesis (Secuencia 1)
	tenantID := "tenant-alpha"
	prevHash0 := "0000000000000000000000000000000000000000000000000000000000000000"
	payloadHash1 := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	eventID1 := "evt-1"
	type1 := "ExpedienteCreado"
	time1 := "2026-09-06T10:00:00Z"

	preimage1 := fmt.Sprintf("%s%s%s%s%s", prevHash0, eventID1, type1, payloadHash1, time1)
	hash1Bytes := sha256.Sum256([]byte(preimage1))
	chainHash1 := hex.EncodeToString(hash1Bytes[:])

	_, err = db.Exec(`INSERT INTO audit_hash_chain (event_id, tenant_id, node_id, actor_id, event_type, canonical_payload_hash, previous_hash, chain_hash, created_at_utc)
		VALUES (?, ?, 'node-1', 'user-1', ?, ?, ?, ?, ?)`,
		eventID1, tenantID, type1, payloadHash1, prevHash0, chainHash1, time1)
	if err != nil {
		t.Fatalf("Fallo inserción evento 1: %v", err)
	}

	// Insertar Evento Secundario (Secuencia 2)
	payloadHash2 := "f4c1d2e389fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	eventID2 := "evt-2"
	type2 := "CalculoEjecutado"
	time2 := "2026-09-06T10:05:00Z"

	preimage2 := fmt.Sprintf("%s%s%s%s%s", chainHash1, eventID2, type2, payloadHash2, time2)
	hash2Bytes := sha256.Sum256([]byte(preimage2))
	chainHash2 := hex.EncodeToString(hash2Bytes[:])

	_, err = db.Exec(`INSERT INTO audit_hash_chain (event_id, tenant_id, node_id, actor_id, event_type, canonical_payload_hash, previous_hash, chain_hash, created_at_utc)
		VALUES (?, ?, 'node-1', 'user-1', ?, ?, ?, ?, ?)`,
		eventID2, tenantID, type2, payloadHash2, chainHash1, chainHash2, time2)
	if err != nil {
		t.Fatalf("Fallo inserción evento 2: %v", err)
	}

	// ASERCIÓN 1: El Trigger de inmutabilidad debe rechazar sentencias UPDATE directas
	_, err = db.Exec(`UPDATE audit_hash_chain SET previous_hash = 'corrupted' WHERE sequence_id = 1`)
	if err == nil {
		t.Fatalf("FALLO CRÍTICO: La base de datos permitió UPDATE en audit_hash_chain")
	}

	// ASERCIÓN 2: Verificación algorítmica de la cadena completa (Re-calculo)
	rows, err := db.Query(`SELECT event_id, event_type, canonical_payload_hash, previous_hash, chain_hash, created_at_utc FROM audit_hash_chain ORDER BY sequence_id ASC`)
	if err != nil {
		t.Fatalf("Error leyendo cadena: %v", err)
	}
	defer rows.Close()

	expectedPrev := prevHash0
	for rows.Next() {
		var evtID, evtType, payloadHash, prevHash, chainHash, createdAt string
		if err := rows.Scan(&evtID, &evtType, &payloadHash, &prevHash, &chainHash, &createdAt); err != nil {
			t.Fatalf("Scan error: %v", err)
		}

		if prevHash != expectedPrev {
			t.Fatalf("INCONSISTENCIA DE CADENA: Esperado prev_hash %s, obtenido %s", expectedPrev, prevHash)
		}

		checkPreimage := fmt.Sprintf("%s%s%s%s%s", prevHash, evtID, evtType, payloadHash, createdAt)
		checkHashArr := sha256.Sum256([]byte(checkPreimage))
		checkHashStr := hex.EncodeToString(checkHashArr[:])

		if checkHashStr != chainHash {
			t.Fatalf("CORRUPCIÓN CRIPTOGRÁFICA: Hash recalculado %s no coincide con guardado %s", checkHashStr, chainHash)
		}

		expectedPrev = chainHash
	}
}
