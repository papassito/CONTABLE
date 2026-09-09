package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // Registra el driver de SQLite libre de CGO
)

// InitDB abre la base de datos SQLite, configura los PRAGMAs de integridad y ejecuta las migraciones iniciales.
func InitDB(dbPath string) (*sql.DB, error) {
	// 1. Conectar a la base de datos usando el driver "sqlite"
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la base de datos SQLite: %w", err)
	}

	// Configurar límites del pool de conexiones para asegurar la estabilidad monohilo de escritura de SQLite
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// 2. Habilitar PRAGMAs para asegurar velocidad e integridad contable
	pragmas := []string{
		"PRAGMA foreign_keys = ON;",  // Exige integridad referencial
		"PRAGMA journal_mode = WAL;", // Mejora el rendimiento de lectura/escritura concurrente
		"PRAGMA synchronous = NORMAL;",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			db.Close()
			return nil, fmt.Errorf("error al ejecutar pragma '%s': %w", pragma, err)
		}
	}

	// 3. Crear las tablas principales si no existen
	if err := createSchemas(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("error al crear el esquema de tablas: %w", err)
	}

	return db, nil
}

// createSchemas define la estructura DDL de la base de datos contable.
func createSchemas(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tenant_tenants (
		id TEXT PRIMARY KEY,
		code TEXT NOT NULL UNIQUE,
		legal_name TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'ACTIVE',
		created_at_utc DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS tenant_taxpayers (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL REFERENCES tenant_tenants(id) ON DELETE RESTRICT,
		rfc TEXT NOT NULL,
		legal_name TEXT NOT NULL,
		tax_person_type TEXT NOT NULL,
		state_code TEXT NOT NULL,
		is_active INTEGER NOT NULL DEFAULT 1,
		created_at_utc DATETIME NOT NULL,
		UNIQUE(tenant_id, id),
		UNIQUE(tenant_id, rfc)
	);

	CREATE TABLE IF NOT EXISTS identity_users (
		id TEXT PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		is_superadmin BOOLEAN NOT NULL DEFAULT 0,
		is_active BOOLEAN NOT NULL DEFAULT 1,
		created_at_utc DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS identity_roles (
		id TEXT PRIMARY KEY,
		code TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		description TEXT
	);

	CREATE TABLE IF NOT EXISTS identity_tenant_memberships (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL REFERENCES tenant_tenants(id) ON DELETE RESTRICT,
		user_id TEXT NOT NULL REFERENCES identity_users(id) ON DELETE RESTRICT,
		role_id TEXT NOT NULL REFERENCES identity_roles(id) ON DELETE RESTRICT,
		is_active BOOLEAN NOT NULL DEFAULT 1,
		created_at_utc DATETIME NOT NULL,
		UNIQUE(tenant_id, user_id)
	);

	CREATE TABLE IF NOT EXISTS identity_sessions (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL REFERENCES identity_users(id) ON DELETE CASCADE,
		token_hash TEXT NOT NULL UNIQUE,
		ip_address TEXT NOT NULL,
		user_agent TEXT NOT NULL,
		expires_at_utc DATETIME NOT NULL,
		created_at_utc DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS audit_events (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL REFERENCES tenant_tenants(id) ON DELETE RESTRICT,
		event_type TEXT NOT NULL,
		payload_json TEXT NOT NULL,
		previous_hash TEXT NOT NULL,
		current_hash TEXT NOT NULL,
		created_at_utc DATETIME NOT NULL,
		UNIQUE(tenant_id, previous_hash)
	);

	CREATE TABLE IF NOT EXISTS accounts (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL REFERENCES tenant_tenants(id) ON DELETE RESTRICT,
		code TEXT NOT NULL,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'ACTIVA',
		accepts_move BOOLEAN NOT NULL DEFAULT 1,
		current_balance_cents INTEGER NOT NULL DEFAULT 0,
		created_at_utc DATETIME NOT NULL,
		UNIQUE(tenant_id, code)
	);

	CREATE TABLE IF NOT EXISTS journal_entries (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL REFERENCES tenant_tenants(id) ON DELETE RESTRICT,
		number TEXT NOT NULL,
		date DATETIME NOT NULL,
		concept TEXT NOT NULL,
		reference TEXT NOT NULL,
		status TEXT NOT NULL CHECK (status IN ('BORRADOR', 'CONTABILIZADO', 'ANULADO')),
		created_at_utc DATETIME NOT NULL,
		UNIQUE(tenant_id, number)
	);

	CREATE TABLE IF NOT EXISTS journal_lines (
		id TEXT PRIMARY KEY,
		entry_id TEXT NOT NULL,
		account_id TEXT NOT NULL,
		description TEXT NOT NULL,
		debit_cents INTEGER NOT NULL DEFAULT 0,
		credit_cents INTEGER NOT NULL DEFAULT 0,
		third_party_id TEXT,
		FOREIGN KEY(entry_id) REFERENCES journal_entries(id) ON DELETE RESTRICT,
		FOREIGN KEY(account_id) REFERENCES accounts(id) ON DELETE RESTRICT
	);

	CREATE TABLE IF NOT EXISTS ledger_entries (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL REFERENCES tenant_tenants(id) ON DELETE RESTRICT,
		account_id TEXT NOT NULL,
		entry_id TEXT NOT NULL,
		line_id TEXT NOT NULL,
		amount_cents INTEGER NOT NULL,
		balance_cents INTEGER NOT NULL,
		created_at_utc DATETIME NOT NULL,
		FOREIGN KEY(account_id) REFERENCES accounts(id) ON DELETE RESTRICT,
		FOREIGN KEY(entry_id) REFERENCES journal_entries(id) ON DELETE RESTRICT,
		FOREIGN KEY(line_id) REFERENCES journal_lines(id) ON DELETE RESTRICT
	);

	CREATE TRIGGER IF NOT EXISTS trg_prevent_audit_events_update BEFORE UPDATE ON audit_events 
	BEGIN 
		SELECT RAISE(ABORT, 'VIOLACIÓN CONSTITUCIONAL FCOS v2.2: La cadena de auditoría es inmutable.'); 
	END;

	CREATE TRIGGER IF NOT EXISTS trg_prevent_audit_events_delete BEFORE DELETE ON audit_events 
	BEGIN 
		SELECT RAISE(ABORT, 'VIOLACIÓN CONSTITUCIONAL FCOS v2.2: Prohibido eliminar eventos de auditoría.'); 
	END;

	CREATE TRIGGER IF NOT EXISTS trg_prevent_posted_journal_update BEFORE UPDATE ON journal_entries
	WHEN OLD.status = 'CONTABILIZADO'
	BEGIN
		SELECT RAISE(ABORT, 'VIOLACIÓN CONSTITUCIONAL FCOS v2.2: No se puede modificar un asiento contable ya contabilizado.');
	END;

	CREATE TRIGGER IF NOT EXISTS trg_prevent_posted_journal_delete BEFORE DELETE ON journal_entries
	WHEN OLD.status = 'CONTABILIZADO'
	BEGIN
		SELECT RAISE(ABORT, 'VIOLACIÓN CONSTITUCIONAL FCOS v2.2: No se puede eliminar un asiento contable ya contabilizado.');
	END;

	CREATE TRIGGER IF NOT EXISTS trg_prevent_posted_lines_update BEFORE UPDATE ON journal_lines
	BEGIN
		SELECT RAISE(ABORT, 'VIOLACIÓN CONSTITUCIONAL FCOS v2.2: No se pueden modificar líneas de un asiento ya contabilizado.')
		WHERE (SELECT status FROM journal_entries WHERE id = OLD.entry_id) = 'CONTABILIZADO';
	END;

	CREATE TRIGGER IF NOT EXISTS trg_prevent_posted_lines_delete BEFORE DELETE ON journal_lines
	BEGIN
		SELECT RAISE(ABORT, 'VIOLACIÓN CONSTITUCIONAL FCOS v2.2: No se pueden eliminar líneas de un asiento ya contabilizado.')
		WHERE (SELECT status FROM journal_entries WHERE id = OLD.entry_id) = 'CONTABILIZADO';
	END;

	CREATE TRIGGER IF NOT EXISTS trg_validate_journal_line_tenant_insert BEFORE INSERT ON journal_lines
	BEGIN
		SELECT RAISE(ABORT, 'VIOLACIÓN DE AISLAMIENTO FCOS v2.2: No se puede asociar una cuenta de otro Tenant.')
		WHERE (SELECT tenant_id FROM journal_entries WHERE id = NEW.entry_id) IS NOT 
		      (SELECT tenant_id FROM accounts WHERE id = NEW.account_id);
	END;

	CREATE TRIGGER IF NOT EXISTS trg_validate_journal_line_tenant_update BEFORE UPDATE ON journal_lines
	BEGIN
		SELECT RAISE(ABORT, 'VIOLACIÓN DE AISLAMIENTO FCOS v2.2: No se puede asociar una cuenta de otro Tenant.')
		WHERE (SELECT tenant_id FROM journal_entries WHERE id = NEW.entry_id) IS NOT 
		      (SELECT tenant_id FROM accounts WHERE id = NEW.account_id);
	END;

	CREATE TRIGGER IF NOT EXISTS trg_validate_ledger_entry_tenant_insert BEFORE INSERT ON ledger_entries
	BEGIN
		SELECT RAISE(ABORT, 'VIOLACIÓN DE AISLAMIENTO FCOS v2.2: Discrepancia de Tenant en libro mayor.')
		WHERE NEW.tenant_id IS NOT (SELECT tenant_id FROM accounts WHERE id = NEW.account_id) OR
		      NEW.tenant_id IS NOT (SELECT tenant_id FROM journal_entries WHERE id = NEW.entry_id);
	END;

	CREATE TRIGGER IF NOT EXISTS trg_validate_ledger_entry_tenant_update BEFORE UPDATE ON ledger_entries
	BEGIN
		SELECT RAISE(ABORT, 'VIOLACIÓN DE AISLAMIENTO FCOS v2.2: Discrepancia de Tenant en libro mayor.')
		WHERE NEW.tenant_id IS NOT (SELECT tenant_id FROM accounts WHERE id = NEW.account_id) OR
		      NEW.tenant_id IS NOT (SELECT tenant_id FROM journal_entries WHERE id = NEW.entry_id);
	END;

	CREATE TRIGGER IF NOT EXISTS trg_prevent_ledger_entries_update BEFORE UPDATE ON ledger_entries
	BEGIN
		SELECT RAISE(ABORT, 'VIOLACIÓN CONSTITUCIONAL FCOS v2.2: El Libro Mayor (ledger_entries) es estrictamente inmutable.');
	END;

	CREATE TRIGGER IF NOT EXISTS trg_prevent_ledger_entries_delete BEFORE DELETE ON ledger_entries
	BEGIN
		SELECT RAISE(ABORT, 'VIOLACIÓN CONSTITUCIONAL FCOS v2.2: Prohibido eliminar registros del Libro Mayor (ledger_entries).');
	END;

	-- Índices de Alto Rendimiento para evitar barridos de tabla (Full Table Scans)
	CREATE INDEX IF NOT EXISTS idx_journal_lines_entry_id ON journal_lines(entry_id);
	CREATE INDEX IF NOT EXISTS idx_journal_lines_account_id ON journal_lines(account_id);
	CREATE INDEX IF NOT EXISTS idx_ledger_entries_tenant_account ON ledger_entries(tenant_id, account_id);
	CREATE INDEX IF NOT EXISTS idx_ledger_entries_entry_id ON ledger_entries(entry_id);
	CREATE INDEX IF NOT EXISTS idx_audit_events_latest ON audit_events(tenant_id, created_at_utc DESC);
	`

	_, err := db.Exec(query)
	return err
}
