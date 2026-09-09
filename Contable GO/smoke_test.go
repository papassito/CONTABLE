package main_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/klik/fcos-kernel/internal/domain"
	"github.com/klik/fcos-kernel/internal/repository"
	"github.com/klik/fcos-kernel/internal/service"
	"github.com/klik/fcos-kernel/pkg/database"
	_ "modernc.org/sqlite"
)

func TestMVPEndToEndSmoke(t *testing.T) {
	// 1. Inicializar Base de Datos Real en Memoria (con Pragmas, Triggers e Índices reales)
	db, err := database.InitDB("file:smoke_test.db?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("Fallo al inicializar base de datos de humo: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// 2. Preparar el Tenant de Producción
	tenantID := "default-tenant"
	_, err = db.ExecContext(ctx, `
		INSERT INTO tenant_tenants (id, code, legal_name, status, created_at_utc)
		VALUES (?, 'TENANT_MVP', 'Klik Contable MVP Tenant', 'ACTIVE', ?)`,
		tenantID, time.Now().UTC())
	if err != nil {
		t.Fatalf("Error al crear Tenant real: %v", err)
	}

	// 3. Inicializar Repositorios y Servicios Reales (0% Simulación)
	accountRepo := repository.NewAccountRepository(db)
	journalRepo := repository.NewJournalRepository(db)
	ledgerRepo := repository.NewLedgerRepository(db)
	uow := database.NewUnitOfWork(db)
	auditSvc := service.NewDBAuditService(db)

	accountSvc := service.NewAccountService(accountRepo)
	journalSvc := service.NewJournalService(journalRepo, accountRepo, ledgerRepo, uow, auditSvc)

	// 4. Crear Cuentas Auxiliares Activas para el flujo de Partida Doble
	// Cuenta 1: Caja (Activo)
	accCajaID := uuid.New().String()
	accCaja := &domain.Account{
		ID:          accCajaID,
		Code:        "11050501",
		Name:        "Caja General MVP",
		Type:        domain.AccountTypeActivo,
		Status:      domain.AccountStatusActiva,
		AcceptsMove: true,
		CurrentBal:  0,
		CreatedAt:   time.Now(),
	}
	if err := accountSvc.CreateAccount(ctx, accCaja); err != nil {
		t.Fatalf("Fallo al crear cuenta Caja: %v", err)
	}

	// Cuenta 2: Capital Social (Patrimonio)
	accCapitalID := uuid.New().String()
	accCapital := &domain.Account{
		ID:          accCapitalID,
		Code:        "31050501",
		Name:        "Capital Suscrito MVP",
		Type:        domain.AccountTypePatrimonio,
		Status:      domain.AccountStatusActiva,
		AcceptsMove: true,
		CurrentBal:  0,
		CreatedAt:   time.Now(),
	}
	if err := accountSvc.CreateAccount(ctx, accCapital); err != nil {
		t.Fatalf("Fallo al crear cuenta Capital: %v", err)
	}

	// 5. Crear un Asiento en Borrador (Draft) equilibrado
	entryID := uuid.New().String()
	entry := &domain.JournalEntry{
		ID:        entryID,
		Number:    "ASE-0001",
		Date:      time.Now(),
		Concept:   "Aportación Inicial de Socios MVP",
		Reference: "SOCIOS-01",
		Lines: []domain.JournalLine{
			{
				ID:          uuid.New().String(),
				AccountID:   accCajaID,
				Description: "Ingreso en Caja",
				Debit:       5000000, // $50,000.00 pesos (en centavos)
				Credit:      0,
			},
			{
				ID:          uuid.New().String(),
				AccountID:   accCapitalID,
				Description: "Capitalización Inicial",
				Debit:       0,
				Credit:      5000000, // $50,000.00 pesos (en centavos)
			},
		},
	}

	draft, err := journalSvc.CreateDraft(ctx, entry)
	if err != nil {
		t.Fatalf("Fallo al crear Borrador del asiento: %v", err)
	}

	if draft.Status != domain.EntryStatusBorrador {
		t.Fatalf("Se esperaba estado BORRADOR, obtenido: %s", draft.Status)
	}

	// 6. Contabilizar / Postear el Asiento
	err = journalSvc.PostEntry(ctx, draft.ID)
	if err != nil {
		t.Fatalf("Fallo al contabilizar el asiento contable: %v", err)
	}

	// 7. Aserciones del Estado de la Base de Datos post-Posteo
	// A) Verificar estados de saldos
	cajaPost, _ := accountSvc.GetAccount(ctx, accCajaID)
	capitalPost, _ := accountSvc.GetAccount(ctx, accCapitalID)

	if cajaPost.CurrentBal != 5000000 {
		t.Errorf("Balance incorrecto en Caja. Esperado: 5000000, Obtenido: %d", cajaPost.CurrentBal)
	}
	if capitalPost.CurrentBal != 5000000 {
		t.Errorf("Balance incorrecto en Capital. Esperado: 5000000, Obtenido: %d", capitalPost.CurrentBal)
	}

	// B) Verificar integridad del audit log en audit_events (0% Simulación, Criptografía Activa)
	var count int
	var payload string
	var currentHash string
	err = db.QueryRowContext(ctx, "SELECT COUNT(*), payload_json, current_hash FROM audit_events WHERE tenant_id = ?", tenantID).Scan(&count, &payload, &currentHash)
	if err != nil {
		t.Fatalf("Error al consultar la tabla de auditoría real: %v", err)
	}

	if count != 1 {
		t.Errorf("Se esperaba exactamente 1 evento de auditoría real, obtenido: %d", count)
	}

	if currentHash == "" {
		t.Errorf("El hash criptográfico actual no debió generarse en blanco")
	}
}
