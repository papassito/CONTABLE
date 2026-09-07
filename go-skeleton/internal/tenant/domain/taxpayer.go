package domain

import (
	"errors"
	"time"
)

// Taxpayer modela al contribuyente fiscal asociado a un Tenant específico.
type Taxpayer struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	RFC           string    `json:"rfc"`
	LegalName     string    `json:"legal_name"`
	TaxPersonType string    `json:"tax_person_type"` // PHYSICAL, MORAL
	StateCode     string    `json:"state_code"`
	IsActive      bool      `json:"is_active"`
	CreatedAtUTC  time.Time `json:"created_at_utc"`
}

// Validate valida la coherencia mínima del contribuyente.
func (t *Taxpayer) Validate() error {
	if t.ID == "" {
		return errors.New("el ID del contribuyente es requerido")
	}
	if t.TenantID == "" {
		return errors.New("el ID del Tenant es requerido para mantener el aislamiento")
	}
	return nil
}
