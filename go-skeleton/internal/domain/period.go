package domain

import "time"

// CompliancePeriod modela los rangos periódicos fiscales
type CompliancePeriod struct {
	ID           string    `json:"id"`
	PeriodYear   int       `json:"period_year"`
	PeriodMonth  int       `json:"period_month"`
	PeriodType   string    `json:"period_type"`
	PeriodNumber int       `json:"period_number"`
	CreatedAtUTC time.Time `json:"created_at_utc"`
}
