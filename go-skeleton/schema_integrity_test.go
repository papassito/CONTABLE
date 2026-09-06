package schema_test

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Error al abrir base de datos SQLite en memoria: %v", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		t.Fatalf("Error al activar foreign_keys: %v", err)
	}

	schema := `
	CREATE TABLE tenant_tenants (
		id TEXT PRIMARY KEY, code TEXT NOT NULL UNIQUE, legal_name TEXT NOT NULL, status TEXT NOT NULL DEFAULT 'ACTIVE', created_at_utc TEXT NOT NULL DEFAULT '2026-01-01T00:00:00Z'
	);
	CREATE TABLE tenant_taxpayers (
		id TEXT PRIMARY KEY, tenant_id TEXT NOT NULL REFERENCES tenant_tenants(id), rfc TEXT NOT NULL, legal_name TEXT NOT NULL, tax_person_type TEXT NOT NULL, state_code TEXT NOT NULL, is_active INTEGER NOT NULL DEFAULT 1, created_at_utc TEXT NOT NULL DEFAULT '2026-01-01T00:00:00Z', UNIQUE(tenant_id, id), UNIQUE(tenant_id, rfc)
	);
	CREATE TABLE compliance_periods (
		id TEXT PRIMARY KEY, period_year INTEGER NOT NULL, period_month INTEGER NOT NULL, period_type TEXT NOT NULL DEFAULT 'MONTHLY', period_number INTEGER NOT NULL, created_at_utc TEXT NOT NULL DEFAULT '2026-01-01T00:00:00Z', UNIQUE(period_year, period_type, period_number)
	);
	CREATE TABLE compliance_expedientes (
		id TEXT PRIMARY KEY, tenant_id TEXT NOT NULL, taxpayer_id TEXT NOT NULL, period_id TEXT NOT NULL REFERENCES compliance_periods(id), authority_domain TEXT NOT NULL, expediente_code TEXT NOT NULL, current_stage TEXT NOT NULL DEFAULT 'CREATED', is_closed INTEGER NOT NULL DEFAULT 0, created_at_utc TEXT NOT NULL DEFAULT '2026-01-01T00:00:00Z', updated_at_utc TEXT NOT NULL DEFAULT '2026-01-01T00:00:00Z', UNIQUE(tenant_id, id), FOREIGN KEY(tenant_id, taxpayer_id) REFERENCES tenant_taxpayers(tenant_id, id), UNIQUE(tenant_id, taxpayer_id, period_id, authority_domain)
	);
	CREATE TABLE secret_envelopes (
		id TEXT PRIMARY KEY, tenant_id TEXT NOT NULL REFERENCES tenant_tenants(id), secret_type TEXT NOT NULL, encrypted_payload BLOB NOT NULL, wrapped_tek_node BLOB NOT NULL, wrapped_tek_recovery BLOB, payload_nonce BLOB NOT NULL, key_version TEXT NOT NULL, created_at_utc TEXT NOT NULL DEFAULT '2026-01-01T00:00:00Z', updated_at_utc TEXT NOT NULL DEFAULT '2026-01-01T00:00:00Z', UNIQUE(tenant_id, secret_type)
	);
	CREATE TABLE audit_hash_chain (
		sequence_id INTEGER PRIMARY KEY AUTOINCREMENT, event_id TEXT NOT NULL UNIQUE, tenant_id TEXT NOT NULL REFERENCES tenant_tenants(id), node_id TEXT NOT NULL, actor_id TEXT NOT NULL, event_type TEXT NOT NULL, canonical_payload_hash TEXT NOT NULL, previous_hash TEXT NOT NULL, chain_hash TEXT NOT NULL, created_at_utc TEXT NOT NULL DEFAULT '2026-01-01T00:00:00Z'
	);
	CREATE TRIGGER trg_prevent_audit_update BEFORE UPDATE ON audit_hash_chain BEGIN SELECT RAISE(ABORT, 'VIOLACIÓN CONSTITUCIONAL FCOS v2.2: La cadena de auditoría es inmutable.'); END;
	CREATE TRIGGER trg_prevent_audit_delete BEFORE DELETE ON audit_hash_chain BEGIN SELECT RAISE(ABORT, 'VIOLACIÓN CONSTITUCIONAL FCOS v2.2: Prohibido eliminar eventos de auditoría.'); END;
	`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Error al inicializar esquema DDL: %v", err)
	}

	return db
}

// 1. Test de Inmutabilidad en la Cadena de Auditoría
func TestAuditLedgerImmutability(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Inserción válida de Tenant y Evento
	_, err := db.Exec(`INSERT INTO tenant_tenants (id, code, legal_name) VALUES ('t1', 'TENANT_A', 'Empresa A')`)
	if err != nil {
		t.Fatalf("Error al insertar tenant: %v", err)
	}

	_, err = db.Exec(`INSERT INTO audit_hash_chain (event_id, tenant_id, node_id, actor_id, event_type, canonical_payload_hash, previous_hash, chain_hash) 
		VALUES ('e1', 't1', 'node-1', 'user-1', 'CALCULATION_EXECUTED', 'hash_payload', 'hash_prev', 'hash_chain')`)
	if err != nil {
		t.Fatalf("Error al insertar evento de auditoría: %v", err)
	}

	// Aserción 1: Intento de UPDATE debe ser rechazado por el Trigger
	_, err = db.Exec(`UPDATE audit_hash_chain SET previous_hash = 'tampered_hash' WHERE event_id = 'e1'`)
	if err == nil {
		t.Fatalf("FALLO DE SEGURIDAD: Se permitió modificar un registro inmutable en audit_hash_chain")
	}

	// Aserción 2: Intento de DELETE debe ser rechazado por el Trigger
	_, err = db.Exec(`DELETE FROM audit_hash_chain WHERE event_id = 'e1'`)
	if err == nil {
		t.Fatalf("FALLO DE SEGURIDAD: Se permitió eliminar un registro inmutable en audit_hash_chain")
	}
}

// 2. Test de Aislamiento Cruzado entre Tenants (Cross-Tenant Isolation)
func TestCrossTenantIsolationConstraint(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Crear Tenant A y Tenant B
	db.Exec(`INSERT INTO tenant_tenants (id, code, legal_name) VALUES ('tenant_a', 'TA', 'Tenant A'), ('tenant_b', 'TB', 'Tenant B')`)
	db.Exec(`INSERT INTO compliance_periods (id, period_year, period_month, period_type, period_number) VALUES ('p1', 2026, 1, 'MONTHLY', 1)`)

	// Taxpayer B pertenece exclusivamente a Tenant B
	_, err := db.Exec(`INSERT INTO tenant_taxpayers (id, tenant_id, rfc, legal_name, tax_person_type, state_code) 
		VALUES ('taxpayer_b', 'tenant_b', 'XAXX010101000', 'Taxpayer B', 'MORAL', 'SON')`)
	if err != nil {
		t.Fatalf("Error al insertar taxpayer: %v", err)
	}

	// Aserción: Intentar crear un Expediente asignado a Tenant A pero referenciando a Taxpayer B
	_, err = db.Exec(`INSERT INTO compliance_expedientes (id, tenant_id, taxpayer_id, period_id, authority_domain, expediente_code) 
		VALUES ('exp_1', 'tenant_a', 'taxpayer_b', 'p1', 'SAT', 'EXP-001')`)

	if err == nil {
		t.Fatalf("FALLO DE ISOLAMIENTO: Se permitió asociar un Taxpayer de Tenant B en un Expediente de Tenant A")
	}
}

// 3. Test de Unicidad de Secretos Fiscales por Tenant
func TestSecretsVaultUniqueness(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	db.Exec(`INSERT INTO tenant_tenants (id, code, legal_name) VALUES ('tenant_a', 'TA', 'Tenant A')`)

	// Inserción del primer secreto tipo CIEC
	_, err := db.Exec(`INSERT INTO secret_envelopes (id, tenant_id, secret_type, encrypted_payload, wrapped_tek_node, payload_nonce, key_version) 
		VALUES ('s1', 'tenant_a', 'CIEC', 'data1', 'tek1', 'nonce1', 'v1')`)
	if err != nil {
		t.Fatalf("Error al insertar secreto inicial: %v", err)
	}

	// Aserción: Intentar duplicar secreto tipo CIEC para el mismo Tenant
	_, err = db.Exec(`INSERT INTO secret_envelopes (id, tenant_id, secret_type, encrypted_payload, wrapped_tek_node, payload_nonce, key_version) 
		VALUES ('s2', 'tenant_a', 'CIEC', 'data2', 'tek2', 'nonce2', 'v1')`)

	if err == nil {
		t.Fatalf("FALLO DE BÓVEDA: Se permitió registrar duplicados de secret_type 'CIEC' para el mismo Tenant")
	}
}
