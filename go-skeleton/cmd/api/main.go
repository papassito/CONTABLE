package main

import (
	"database/sql"
	"log"
	"net/http"

	httphandler "github.com/klik/fcos-kernel/internal/handler/http"
	"github.com/klik/fcos-kernel/internal/repository"
	"github.com/klik/fcos-kernel/internal/service"
	"github.com/klik/fcos-kernel/pkg/database"
	_ "modernc.org/sqlite"
)

func main() {
	// 1. Initialize SQLite database with WAL (Write-Ahead Logging) mode
	db, err := sql.Open("sqlite", "file:fcos_local.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		log.Fatalf("Critical failure connecting to SQLite database: %v", err)
	}
	defer db.Close()

	// 2. Instantiate Concrete Persistence Repositories
	accountRepo := repository.NewAccountRepository(db)
	journalRepo := repository.NewJournalRepository(db)
	ledgerRepo := repository.NewLedgerRepository(db)
	uow := database.NewUnitOfWork(db)

	// 3. Inject real dependencies into Application Services
	auditService := service.NewLoggingAuditService()
	accountSvc := service.NewAccountService(accountRepo)
	journalSvc := service.NewJournalService(journalRepo, accountRepo, ledgerRepo, uow, auditService)

	// 4. Configure API Router and Handlers
	router := httphandler.NewRouter(accountSvc, journalSvc)
	archHandler := httphandler.NewArchitectureHandler("./")

	// Endpoint for dynamic inspection from the frontend
	http.HandleFunc("/api/v1/architecture/files", archHandler.GetFileTree)
	http.Handle("/", router)

	log.Println("🚀 FCOS v2.2 Kernel started successfully on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Critical error in HTTP server: %v", err)
	}
}
