package domain

import "time"

// Expediente unifica la carga de periodos y su ciclo de vida operacional
type Expediente struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenant_id"`
	TaxpayerID      string    `json:"taxpayer_id"`
	PeriodID        string    `json:"period_id"`
	AuthorityDomain string    `json:"authority_domain"`
	ExpedienteCode  string    `json:"expediente_code"`
	CurrentStage    string    `json:"current_stage"`
	IsClosed        bool      `json:"is_closed"`
	CreatedAtUTC    time.Time `json:"created_at_utc"`
	UpdatedAtUTC    time.Time `json:"updated_at_utc"`
}
