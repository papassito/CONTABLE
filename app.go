package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/klik/fcos-kernel/internal/domain"
	"github.com/klik/fcos-kernel/internal/repository"
	"github.com/klik/fcos-kernel/internal/service"
	"github.com/klik/fcos-kernel/pkg/database"
	_ "modernc.org/sqlite"
)

// App representa la estructura principal de la aplicación Wails.
type App struct {
	ctx        context.Context
	ctxMu      sync.RWMutex
	db         *sql.DB
	accountSvc service.AccountService
	journalSvc service.JournalService
}

// NewApp inicializa la base de datos SQLite, los repositorios y los servicios del dominio.
func NewApp() *App {
	// 1. Inicializar la base de datos SQLite local
	dbDir, err := os.UserConfigDir()
	if err != nil {
		dbDir = "."
	}
	appDir := filepath.Join(dbDir, "ContableFix")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		log.Fatalf("Fallo crítico al crear el directorio de la base de datos: %v", err)
	}
	dbPath := filepath.Join(appDir, "fcos_local.db")

	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Fallo crítico al conectar con la base de datos SQLite: %v", err)
	}

	// 2. Instanciar repositorios de persistencia
	accountRepo := repository.NewAccountRepository(db)
	journalRepo := repository.NewJournalRepository(db)
	ledgerRepo := repository.NewLedgerRepository(db)
	uow := database.NewUnitOfWork(db)

	// 3. Inyectar servicios del dominio contable
	auditService := service.NewDBAuditService(db)
	accountSvc := service.NewAccountService(accountRepo)
	journalSvc := service.NewJournalService(journalRepo, accountRepo, ledgerRepo, uow, auditService)

	return &App{
		ctx:        context.Background(),
		db:         db,
		accountSvc: accountSvc,
		journalSvc: journalSvc,
	}
}

// startup se ejecuta al iniciar la ventana gráfica de Wails.
func (a *App) startup(ctx context.Context) {
	a.ctxMu.Lock()
	defer a.ctxMu.Unlock()
	a.ctx = ctx
}

// shutdown libera las conexiones a la base de datos al cerrar la aplicación.
func (a *App) shutdown(ctx context.Context) {
	a.ctxMu.Lock()
	defer a.ctxMu.Unlock()
	if a.db != nil {
		a.db.Close()
	}
}

// getCtx obtiene un contexto seguro y garantizado, evitando referencias nulas.
func (a *App) getCtx() context.Context {
	a.ctxMu.RLock()
	defer a.ctxMu.RUnlock()
	if a.ctx == nil {
		return context.Background()
	}
	return a.ctx
}

// CreateDraft guarda un borrador de asiento contable en la base de datos real.
func (a *App) CreateDraft(entry *domain.JournalEntry) (*domain.JournalEntry, error) {
	res, err := a.journalSvc.CreateDraft(a.getCtx(), entry)
	return res, sanitizeError(err)
}

// PostEntry contabiliza de forma definitiva e inmutable un asiento.
func (a *App) PostEntry(entryID string) error {
	err := a.journalSvc.PostEntry(a.getCtx(), entryID)
	return sanitizeError(err)
}

// ReverseEntry genera un contraasiento de anulación con registro de auditoría.
func (a *App) ReverseEntry(entryID string, reason string) (*domain.JournalEntry, error) {
	res, err := a.journalSvc.ReverseEntry(a.getCtx(), entryID, reason)
	return res, sanitizeError(err)
}

// GetEntry obtiene los detalles de un asiento contable mediante su ID.
func (a *App) GetEntry(id string) (*domain.JournalEntry, error) {
	res, err := a.journalSvc.GetEntry(a.getCtx(), id)
	return res, sanitizeError(err)
}

// CreateAccount registra una nueva cuenta contable en el catálogo.
func (a *App) CreateAccount(account *domain.Account) error {
	err := a.accountSvc.CreateAccount(a.getCtx(), account)
	return sanitizeError(err)
}

// GetAccount obtiene una cuenta contable mediante su identificador único.
func (a *App) GetAccount(id string) (*domain.Account, error) {
	res, err := a.accountSvc.GetAccount(a.getCtx(), id)
	return res, sanitizeError(err)
}

// GetAccountByCode obtiene una cuenta contable mediante su código numérico (ej: "110505").
func (a *App) GetAccountByCode(code string) (*domain.Account, error) {
	res, err := a.accountSvc.GetAccountByCode(a.getCtx(), code)
	return res, sanitizeError(err)
}

// ListAccounts obtiene la lista completa de cuentas contables registradas.
func (a *App) ListAccounts() ([]*domain.Account, error) {
	res, err := a.accountSvc.ListAccounts(a.getCtx())
	return res, sanitizeError(err)
}

// sanitizeError oculta detalles internos de la base de datos para proteger la seguridad del sistema.
func sanitizeError(err error) error {
	if err == nil {
		return nil
	}
	errMsg := err.Error()

	// Permitir que los mensajes constitucionales lleguen a la interfaz de usuario
	if strings.Contains(strings.ToUpper(errMsg), "VIOLACIÓN") {
		return err
	}

	lowered := strings.ToLower(errMsg)
	isSystemOrDB := strings.Contains(lowered, "sqlite") ||
		strings.Contains(lowered, "database is locked") ||
		strings.Contains(lowered, "no such table") ||
		strings.Contains(lowered, "file:")

	if isSystemOrDB {
		log.Printf("[SECURITY-ERROR-LEAK] Error interno sanitizado: %v", err)
		return fmt.Errorf("error interno del sistema de base de datos")
	}

	return err
}
