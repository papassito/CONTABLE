package domain

import "time"

// Tenant representa a un inquilino multi-tenant aislado
type Tenant struct {
	ID           string    `json:"id"`
	Code         string    `json:"code"`
	LegalName    string    `json:"legal_name"`
	Status       string    `json:"status"`
	CreatedAtUTC time.Time `json:"created_at_utc"`
}
