package http

import (
	"net/http"

	"github.com/klik/fcos-kernel/internal/service"
)

// NewRouter creates the HTTP router and connects handlers with their services.
func NewRouter(accountSvc service.AccountService, journalSvc service.JournalService) http.Handler {
	mux := http.NewServeMux()

	// Create handlers with their injected services
	accountHandler := NewAccountHandler(accountSvc)
	journalHandler := NewJournalHandler(journalSvc)

	// Chart of Accounts
	mux.HandleFunc("POST /api/v1/accounts", accountHandler.Create)
	mux.HandleFunc("GET /api/v1/accounts", accountHandler.List)
	mux.HandleFunc("GET /api/v1/accounts/{id}", accountHandler.GetByID)

	// Journal Entries
	mux.HandleFunc("POST /api/v1/journal-entries", journalHandler.CreateDraft)
	mux.HandleFunc("POST /api/v1/journal-entries/{id}/post", journalHandler.Post)
	mux.HandleFunc("POST /api/v1/journal-entries/{id}/reverse", journalHandler.Reverse)
	mux.HandleFunc("GET /api/v1/journal-entries/{id}", journalHandler.GetByID)

	return mux
}
