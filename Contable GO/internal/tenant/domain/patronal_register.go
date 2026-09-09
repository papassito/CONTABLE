package domain

import (
	"errors"
	"time"
)

// PatronalRegister representa un registro patronal del IMSS asociado a un contribuyente.
type PatronalRegister struct {
	ID                    string    `json:"id"`
	TenantID              string    `json:"tenant_id"`
	TaxpayerID            string    `json:"taxpayer_id"`
	RegistroPatronal      string    `json:"registro_patronal"`
	ImssSubdelegationCode string    `json:"imss_subdelegation_code"`
	RiskClass             string    `json:"risk_class"`
	RiskPremium           float64   `json:"risk_premium"`
	IsActive              bool      `json:"is_active"`
	CreatedAtUTC          time.Time `json:"created_at_utc"`
}

// Validate verifica la validez estructural de un registro patronal.
func (p *PatronalRegister) Validate() error {
	if p.ID == "" {
		return errors.New("el ID del registro patronal es requerido")
	}
	if p.TenantID == "" || p.TaxpayerID == "" {
		return errors.New("los identificadores de Tenant y Taxpayer son requeridos")
	}
	return nil
}
