package domain

import "time"

// ParameterVersion modela parámetros oficiales (e.g. UMA) con soporte de inmutabilidad
type ParameterVersion struct {
	ID                string     `json:"id"`
	SourceID          string     `json:"source_id"`
	ParameterCode     string     `json:"parameter_code"`
	Jurisdiction      string     `json:"jurisdiction"`
	NumericValue      float64    `json:"numeric_value"`
	EffectiveFrom     time.Time  `json:"effective_from"`
	EffectiveTo       *time.Time `json:"effective_to"`
	PublicationDate   time.Time  `json:"publication_date"`
	DocumentReference string     `json:"document_reference"`
	DocumentHash      string     `json:"document_hash"`
	IsFrozen          bool       `json:"is_frozen"` // Ley de inmutabilidad contable
	CreatedAtUTC      time.Time  `json:"created_at_utc"`
}
