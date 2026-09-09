package http

import (
	"github.com/klik/fcos-kernel/internal/service"
	"net/http"
)

type JournalHandler struct {
	service service.JournalService
}

func NewJournalHandler(service service.JournalService) *JournalHandler {
	return &JournalHandler{service: service}
}

func (h *JournalHandler) CreateDraft(w http.ResponseWriter, r *http.Request) {
	// TODO: Recibir asiento contable y llamar a CreateDraft
}

func (h *JournalHandler) Post(w http.ResponseWriter, r *http.Request) {
	// TODO: Contabilizar asiento y asentar en libro mayor
}

func (h *JournalHandler) Reverse(w http.ResponseWriter, r *http.Request) {
	// TODO: Revertir asiento contable
}

func (h *JournalHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Consultar asiento específico
}
