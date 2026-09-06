package main

import (
	"context"
	"testing"

	"github.com/klik/fcos-kernel/config"
	"github.com/klik/fcos-kernel/pkg/database"
)

func TestEngineBootstrapAndUnitOfWorkExecution(t *testing.T) {
	cfg := &config.Config{
		Environment:    "test",
		DatabaseDriver: "sqlite",
		DatabaseURL:    ":memory:",
		NodeID:         "test-node-01",
	}

	engine, err := BootstrapEngine(cfg)
	if err != nil {
		t.Fatalf("FALLO DE ARRANQUE: No se pudo inicializar el motor FCOS: %v", err)
	}
	defer engine.DB.Close()

	// Crear tabla de prueba dentro de SQLite en memoria
	_, err = engine.DB.Exec(`CREATE TABLE test_isolation (id TEXT PRIMARY KEY, val TEXT);`)
	if err != nil {
		t.Fatalf("Error al crear tabla de prueba: %v", err)
	}

	// Probar ejecución transaccional exitosa vía UnitOfWork
	err = engine.UOW.Execute(context.Background(), func(txCtx context.Context) error {
		tx, ok := database.GetTx(txCtx)
		if !ok {
			t.Fatalf("ERROR DE UOW: Transacción no encontrada en el contexto")
		}

		_, execErr := tx.ExecContext(txCtx, `INSERT INTO test_isolation (id, val) VALUES ('1', 'ok');`)
		return execErr
	})

	if err != nil {
		t.Fatalf("FALLO DE UOW: La transacción no pudo ser confirmada: %v", err)
	}

	// Verificar persistencia post-commit
	var count int
	err = engine.DB.QueryRow(`SELECT COUNT(*) FROM test_isolation WHERE id = '1';`).Scan(&count)
	if err != nil || count != 1 {
		t.Fatalf("FALLO DE PERSISTENCIA: El registro no se confirmó correctamente vía UOW")
	}
}
