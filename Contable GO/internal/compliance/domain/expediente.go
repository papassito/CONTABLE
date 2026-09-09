package domain

import "time"

// Expediente representa la agregación de obligaciones y estatus de cumplimiento para un periodo específico.
type Expediente struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenant_id"`
	TaxpayerID      string    `json:"taxpayer_id"`
	PeriodID        string    `json:"period_id"`
	AuthorityDomain string    `json:"authority_domain"` // SAT, IMSS, ISRTP_ESTATAL
	CurrentStage    string    `json:"current_stage"`
	IsClosed        bool      `json:"is_closed"`
	CreatedAtUTC    time.Time `json:"created_at_utc"`
	UpdatedAtUTC    time.Time `json:"updated_at_utc"`
}

// CanBeModified verifica si el expediente permite mutaciones basadas en su estado de cierre.
func (e *Expediente) CanBeModified() bool {
	return !e.IsClosed
}
