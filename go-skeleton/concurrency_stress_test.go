package fcos_kernel_test

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"

	"github.com/klik/fcos-kernel/pkg/database"
	_ "modernc.org/sqlite"
)

func TestHighConcurrencyACIDAndTenantIsolation(t *testing.T) {
	db, err := sql.Open("sqlite", "file:stress_test.db?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("Error abriendo DB: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec("PRAGMA foreign_keys = ON; PRAGMA journal_mode = WAL;"); err != nil {
		t.Fatalf("Error en pragmas: %v", err)
	}

	schema := `
	CREATE TABLE tenant_balances (
		tenant_id TEXT PRIMARY KEY,
		total_calculated_cents INTEGER NOT NULL DEFAULT 0
	);
	`
	db.Exec(schema)

	tenantA := "tenant-a"
	tenantB := "tenant-b"
	db.Exec(`INSERT INTO tenant_balances (tenant_id, total_calculated_cents) VALUES (?, 0), (?, 0)`, tenantA, tenantB)

	uow := database.NewUnitOfWork(db)
	var wg sync.WaitGroup

	numWorkers := 20
	iterationsPerWorker := 10
	amountPerCalc := int64(100) // 100 céntimos ($1.00)

	// Lanza trabajadores concurrentes para Tenant A y Tenant B
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		targetTenant := tenantA
		if i%2 == 0 {
			targetTenant = tenantB
		}

		go func(workerID int, tID string) {
			defer wg.Done()
			for j := 0; j < iterationsPerWorker; j++ {
				_ = uow.Execute(context.Background(), func(txCtx context.Context) error {
					tx, ok := database.GetTx(txCtx)
					if !ok {
						return fmt.Errorf("Tx no encontrada")
					}

					_, execErr := tx.ExecContext(txCtx, `
						UPDATE tenant_balances 
						SET total_calculated_cents = total_calculated_cents + ? 
						WHERE tenant_id = ?;`,
						amountPerCalc, tID)

					return execErr
				})
			}
		}(i, targetTenant)
	}

	wg.Wait()

	// Aserciones de Cuadre Exacto post-estresamiento
	var balA, balB int64
	db.QueryRow(`SELECT total_calculated_cents FROM tenant_balances WHERE tenant_id = ?`, tenantA).Scan(&balA)
	db.QueryRow(`SELECT total_calculated_cents FROM tenant_balances WHERE tenant_id = ?`, tenantB).Scan(&balB)

	expectedPerTenant := int64(numWorkers/2) * int64(iterationsPerWorker) * amountPerCalc

	if balA != expectedPerTenant {
		t.Fatalf("CARRERA DE DATOS EN TENANT A: Esperado %d, Obtenido %d", expectedPerTenant, balA)
	}

	if balB != expectedPerTenant {
		t.Fatalf("CARRERA DE DATOS EN TENANT B: Esperado %d, Obtenido %d", expectedPerTenant, balB)
	}
}
