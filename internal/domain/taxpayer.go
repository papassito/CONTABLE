package domain

import "time"

// Taxpayer modela al contribuyente fiscal asociado a un Tenant
type Taxpayer struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	RFC           string    `json:"rfc"`
	LegalName     string    `json:"legal_name"`
	TaxPersonType string    `json:"tax_person_type"`
	StateCode     string    `json:"state_code"`
	IsActive      bool      `json:"is_active"`
	CreatedAtUTC  time.Time `json:"created_at_utc"`
}
