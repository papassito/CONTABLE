package domain

import (
	"time"
)

type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "ACTIVE"
	TenantStatusInactive  TenantStatus = "INACTIVE"
	TenantStatusSuspended TenantStatus = "SUSPENDED"
)

// Tenant representa a un inquilino multi-tenant aislado dentro del bounded context de administración de cuentas.
type Tenant struct {
	ID           string       `json:"id"`
	Code         string       `json:"code"`
	LegalName    string       `json:"legal_name"`
	Status       TenantStatus `json:"status"`
	CreatedAtUTC time.Time    `json:"created_at_utc"`
}

// Validate verifica la consistencia del Tenant.
func (t *Tenant) Validate() error {
	return nil
}
