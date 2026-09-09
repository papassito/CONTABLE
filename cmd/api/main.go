package main

import (
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
	db, err := database.InitDB("file:fcos_local.db?_pragma=busy_timeout(5000)")
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
	auditService := service.NewDBAuditService(db)
	accountSvc := service.NewAccountService(accountRepo)
	journalSvc := service.NewJournalService(journalRepo, accountRepo, ledgerRepo, uow, auditService)

	// 4. Configure API Router and Handlers
	router := httphandler.NewRouter(accountSvc, journalSvc)
	archHandler := httphandler.NewArchitectureHandler("./")

	mux := http.NewServeMux()
	// Endpoint for dynamic inspection from the frontend
	mux.HandleFunc("/api/v1/architecture/files", archHandler.GetFileTree)
	mux.Handle("/", router)

	log.Println("🚀 FCOS v2.2 Kernel started successfully on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Critical error in HTTP server: %v", err)
	}
}
