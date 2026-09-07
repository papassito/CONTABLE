package domain

import "time"

// PatronalRegister representa un registro patronal del IMSS
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
