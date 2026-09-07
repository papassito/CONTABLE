package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/klik/fcos-kernel/internal/domain"
	"github.com/klik/fcos-kernel/internal/repository"
	"github.com/klik/fcos-kernel/internal/service"
	"github.com/klik/fcos-kernel/pkg/database"
	_ "modernc.org/sqlite"
)

// App struct represents the Wails application
type App struct {
	ctx        context.Context
	db         *sql.DB
	accountSvc service.AccountService
	journalSvc service.JournalService
}

// NewApp creates a new App application struct
func NewApp() *App {
	// 1. Initialize SQLite database with WAL (Write-Ahead Logging)
	db, err := sql.Open("sqlite", "file:fcos_local.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		log.Fatalf("Critical failure connecting to SQLite database: %v", err)
	}

	// 2. Instantiate Persistence Repositories
	accountRepo := repository.NewAccountRepository(db)
	journalRepo := repository.NewJournalRepository(db)
	ledgerRepo := repository.NewLedgerRepository(db)
	uow := database.NewUnitOfWork(db)

	// 3. Inject Core Domain Services
	auditService := service.NewLoggingAuditService()
	accountSvc := service.NewAccountService(accountRepo)
	journalSvc := service.NewJournalService(journalRepo, accountRepo, ledgerRepo, uow, auditService)

	return &App{
		db:         db,
		accountSvc: accountSvc,
		journalSvc: journalSvc,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown() {
	if a.db != nil {
		a.db.Close()
	}
}

func (a *App) CreateDraft(entry *domain.JournalEntry) (*domain.JournalEntry, error) {
	return a.journalSvc.CreateDraft(a.ctx, entry)
}

func (a *App) PostEntry(entryID string) error {
	return a.journalSvc.PostEntry(a.ctx, entryID)
}

func (a *App) ReverseEntry(entryID string, reason string) (*domain.JournalEntry, error) {
	return a.journalSvc.ReverseEntry(a.ctx, entryID, reason)
}

func (a *App) GetEntry(id string) (*domain.JournalEntry, error) {
	return a.journalSvc.GetEntry(a.ctx, id)
}
