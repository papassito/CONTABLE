package domain

import "time"

// Obligation representa obligaciones normativas asignadas a un Taxpayer
type Obligation struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	TaxpayerID   string    `json:"taxpayer_id"`
	Code         string    `json:"code"`
	Description  string    `json:"description"`
	IsActive     bool      `json:"is_active"`
	DueDateRule  string    `json:"due_date_rule"`
	CreatedAtUTC time.Time `json:"created_at_utc"`
}
