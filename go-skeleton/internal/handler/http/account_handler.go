package http

import (
	"github.com/klik/fcos-kernel/internal/service"
	"net/http"
)

type AccountHandler struct {
	service service.AccountService
}

func NewAccountHandler(service service.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Decodificar payload JSON, validar y llamar a h.service.CreateAccount
}

func (h *AccountHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Extraer ID de la URL y responder JSON
}

func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: Listar cuentas y responder JSON
}
