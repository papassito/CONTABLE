package domain

import "time"

type AccountType string

const (
	AccountTypeActivo     AccountType = "ACTIVO"
	AccountTypePasivo     AccountType = "PASIVO"
	AccountTypePatrimonio AccountType = "PATRIMONIO"
	AccountTypeIngreso    AccountType = "INGRESO"
	AccountTypeGasto      AccountType = "GASTO"
	AccountTypeCosto      AccountType = "COSTO"
)

type AccountStatus string

const (
	AccountStatusActiva    AccountStatus = "ACTIVA"
	AccountStatusInactiva  AccountStatus = "INACTIVA"
	AccountStatusBloqueada AccountStatus = "BLOQUEADA"
)

type Account struct {
	ID          string        `json:"id"`
	Code        string        `json:"code"`
	Name        string        `json:"name"`
	Type        AccountType   `json:"type"`
	ParentID    *string       `json:"parent_id"`
	Level       int           `json:"level"`
	AcceptsMove bool          `json:"accepts_move"`
	Status      AccountStatus `json:"status"`
	CurrentBal  float64       `json:"current_balance"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}
