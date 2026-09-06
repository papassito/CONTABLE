package http

import (
	"net/http"
)

type RouterConfig struct {
	AccountHandler *AccountHandler
	JournalHandler *JournalHandler
}

func NewRouter(cfg RouterConfig) http.Handler {
	mux := http.NewServeMux()

	// Plan de cuentas contable
	mux.HandleFunc("POST /api/v1/accounts", cfg.AccountHandler.Create)
	mux.HandleFunc("GET /api/v1/accounts", cfg.AccountHandler.List)
	mux.HandleFunc("GET /api/v1/accounts/{id}", cfg.AccountHandler.GetByID)

	// Asientos y comprobantes contables
	mux.HandleFunc("POST /api/v1/journal-entries", cfg.JournalHandler.CreateDraft)
	mux.HandleFunc("POST /api/v1/journal-entries/{id}/post", cfg.JournalHandler.Post)
	mux.HandleFunc("POST /api/v1/journal-entries/{id}/reverse", cfg.JournalHandler.Reverse)
	mux.HandleFunc("GET /api/v1/journal-entries/{id}", cfg.JournalHandler.GetByID)

	return mux
}
