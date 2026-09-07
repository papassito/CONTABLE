package main_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	healthApp "github.com/klik/fcos-kernel/internal/health/infrastructure"
	_ "modernc.org/sqlite"
)

func TestDay2SyntheticMonitoringAndHealthCheck(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Error abriendo DB: %v", err)
	}
	defer db.Close()

	schema := `
	CREATE TABLE compliance_expedientes (id TEXT PRIMARY KEY, current_stage TEXT, is_closed INTEGER);
	CREATE TABLE compliance_calculations (id TEXT PRIMARY KEY);
	CREATE TABLE audit_hash_chain (sequence_id INTEGER PRIMARY KEY AUTOINCREMENT);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Error al crear el esquema de pruebas: %v", err)
	}

	// Inserciones sintéticas de simulación
	if _, err := db.Exec(`INSERT INTO compliance_expedientes (id, current_stage, is_closed) VALUES ('e1', 'CREATED', 0), ('e2', 'REQUIRES_HUMAN', 0);`); err != nil {
		t.Fatalf("Error al insertar expedientes sintéticos: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO compliance_calculations (id) VALUES ('c1'), ('c2'), ('c3');`); err != nil {
		t.Fatalf("Error al insertar cálculos sintéticos: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO audit_hash_chain (sequence_id) VALUES (1), (2), (3), (4), (5);`); err != nil {
		t.Fatalf("Error al insertar nodos de auditoría sintéticos: %v", err)
	}

	collector := healthApp.NewTelemetryCollector(db)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	start := time.Now()
	telemetry, err := collector.CollectNodeTelemetry(ctx)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Fallo en la recolección de telemetría Day-2: %v", err)
	}

	// Aserciones de Operatividad y Latencia
	if telemetry.TotalExpedientesActive != 2 {
		t.Fatalf("Métrica de expedientes activos incorrecta: %d", telemetry.TotalExpedientesActive)
	}

	if telemetry.TotalCalculationsCount != 3 {
		t.Fatalf("Métrica de cálculos incorrecta: %d", telemetry.TotalCalculationsCount)
	}

	if telemetry.AuditChainLength != 5 {
		t.Fatalf("Métrica de longitud de cadena incorrecta: %d", telemetry.AuditChainLength)
	}

	if telemetry.UnresolvedRPABarriers != 1 {
		t.Fatalf("Métrica de barreras RPA pendientes incorrecta: %d", telemetry.UnresolvedRPABarriers)
	}

	// Evitar fallos falsos positivos en entornos CI/CD compartidos o lentos
	if elapsed > 100*time.Millisecond && elapsed <= 500*time.Millisecond {
		t.Logf("⚠️ ADVERTENCIA DE RENDIMIENTO: Recolección de telemetría tomó %v (Umbral de advertencia: 100ms)", elapsed)
	} else if elapsed > 500*time.Millisecond {
		t.Fatalf("❌ ALERTA DE RENDIMIENTO CRÍTICO: Recolección de telemetría tomó %v (Límite estricto: 500ms)", elapsed)
	}
}
