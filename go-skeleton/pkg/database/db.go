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
	CREATE TABLE IF NOT EXISTS audit_events (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		event_type TEXT NOT NULL,
		payload_json TEXT NOT NULL,
		previous_hash TEXT NOT NULL,
		current_hash TEXT NOT NULL,
		created_at_utc DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS accounts (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		code TEXT NOT NULL,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		accepts_move BOOLEAN NOT NULL DEFAULT 1,
		current_balance_cents INTEGER NOT NULL DEFAULT 0,
		created_at_utc DATETIME NOT NULL,
		UNIQUE(tenant_id, code)
	);
	`

	_, err := db.Exec(query)
	return err
}
