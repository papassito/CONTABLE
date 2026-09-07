package domain

import (
	"errors"
	"time"
)

// Obligation representa una obligación normativa de cumplimiento (compliance) asignada a un contribuyente.
type Obligation struct {
	ID              string     `json:"id"`
	TenantID        string     `json:"tenant_id"`
	TaxpayerID      string     `json:"taxpayer_id"`
	Code            string     `json:"code"` // E.g., "DIOT", "REP_SUSS", "SIPARE"
	Description     string     `json:"description"`
	AuthorityDomain string     `json:"authority_domain"` // E.g., "SAT", "IMSS"
	DueDate         time.Time  `json:"due_date"`
	IsActive        bool       `json:"is_active"`
	Status          string     `json:"status"` // PENDING, COMPLETED, OVERDUE
	PresentedAtUTC  *time.Time `json:"presented_at_utc,omitempty"`
	CreatedAtUTC    time.Time  `json:"created_at_utc"`
	UpdatedAtUTC    time.Time  `json:"updated_at_utc"`
}

// Validate verifica la consistencia mínima de la obligación normativa.
func (o *Obligation) Validate() error {
	if o.ID == "" {
		return errors.New("el ID de la obligación es requerido")
	}
	if o.TenantID == "" || o.TaxpayerID == "" {
		return errors.New("los identificadores de Tenant y Taxpayer son requeridos para mantener el aislamiento")
	}
	if o.Code == "" {
		return errors.New("el código de la obligación no puede estar vacío")
	}
	return nil
}

// IsOverdue determina si la obligación está vencida basándose en la fecha actual UTC.
func (o *Obligation) IsOverdue() bool {
	if o.Status == "COMPLETED" || o.PresentedAtUTC != nil {
		return false
	}
	return time.Now().UTC().After(o.DueDate)
}
