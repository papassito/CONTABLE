package integration_test

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
	db.Exec(schema)

	// Inserciones sintéticas de simulación
	db.Exec(`INSERT INTO compliance_expedientes (id, current_stage, is_closed) VALUES ('e1', 'CREATED', 0), ('e2', 'REQUIRES_HUMAN', 0);`)
	db.Exec(`INSERT INTO compliance_calculations (id) VALUES ('c1'), ('c2'), ('c3');`)
	db.Exec(`INSERT INTO audit_hash_chain (sequence_id) VALUES (1), (2), (3), (4), (5);`)

	collector := healthApp.NewTelemetryCollector(db)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	start := time.Now()
	telemetry, err := collector.CollectNodeTelemetry()
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

	if elapsed > 100*time.Millisecond {
		t.Fatalf("ALERTA DE RENDIMIENTO DAY-2: Recolección de telemetría tomó %v (Límite: 100ms)", elapsed)
	}
}
