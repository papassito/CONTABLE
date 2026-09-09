package domain

import "errors"

var (
	ErrAccountActiveWithBalance = errors.New("FCOS_ERR_ACCOUNT: No es posible desactivar una cuenta con saldo distinto de cero")
	ErrHierarchyCycle           = errors.New("FCOS_ERR_ACCOUNT: Se detectó un ciclo infinito en la estructura jerárquica de cuentas")
)

type AccountStatus string

const (
	AccountActive   AccountStatus = "ACTIVA"
	AccountInactive AccountStatus = "INACTIVA"
)

// Account define la estructura del catálogo del plan de cuentas en el dominio puro.
type Account struct {
	ID          string        `json:"id"`
	TenantID    string        `json:"tenant_id"`
	Code        string        `json:"code"`
	Name        string        `json:"name"`
	ParentID    string        `json:"parent_id"`
	AcceptsMove bool          `json:"accepts_move"`
	Status      AccountStatus `json:"status"`
	Balance     Cents         `json:"balance_cents"`
}

// ValidateTransition comprueba si la cuenta contable puede transicionar al estado inactivo.
func (a *Account) ValidateTransition(targetStatus AccountStatus) error {
	if targetStatus == AccountInactive && a.Balance != 0 {
		return ErrAccountActiveWithBalance
	}
	return nil
}

// DetectCycle comprueba si un plan de cuentas contiene bucles recursivos perjudiciales.
// Requiere un mapa con el catálogo completo del tenant activo.
func DetectCycle(accounts map[string]*Account, startID string) bool {
	if startID == "" {
		return false
	}

	visited := make(map[string]bool)
	currentID := startID

	for currentID != "" {
		if visited[currentID] {
			return true // Se encontró un bucle de dependencia jerárquica.
		}
		visited[currentID] = true

		acc, exists := accounts[currentID]
		if !exists {
			break // Apunta a una raíz inexistente o finalizada.
		}
		currentID = acc.ParentID
	}
	return false
}
