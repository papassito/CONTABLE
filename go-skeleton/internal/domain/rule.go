package domain

import "time"

// RuleVersion modela las fórmulas de cálculo oficializadas
type RuleVersion struct {
	ID            string     `json:"id"`
	RuleCode      string     `json:"rule_code"`
	VersionLabel  string     `json:"version_label"`
	Description   string     `json:"description"`
	FormulaJSON   string     `json:"formula_json"`
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to"`
	IsFrozen      bool       `json:"is_frozen"`
	CreatedAtUTC  time.Time  `json:"created_at_utc"`
}
