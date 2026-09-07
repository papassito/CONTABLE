package domain

import (
	"errors"
	"time"
)

// CompliancePeriod representa la especificación de un período fiscal en el contexto de cumplimiento.
type CompliancePeriod struct {
	ID           string    `json:"id"`
	PeriodYear   int       `json:"period_year"`
	PeriodMonth  int       `json:"period_month"`
	PeriodType   string    `json:"period_type"` // MONTHLY, BIMONTHLY, ANNUAL, CUSTOM
	IsClosed     bool      `json:"is_closed"`
	CreatedAtUTC time.Time `json:"created_at_utc"`
}

// Validate verifica que el periodo tenga valores cronológicos y fiscales coherentes.
func (p *CompliancePeriod) Validate() error {
	if p.PeriodYear < 2000 {
		return errors.New("el año de cumplimiento no puede ser anterior al 2000")
	}
	return nil
}
