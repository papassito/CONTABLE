package healthApp

import (
	"database/sql"
)

type SystemTelemetry struct {
	TotalExpedientesActive int64   `json:"total_expedientes_active"`
	TotalCalculationsCount int64   `json:"total_calculations_count"`
	AuditChainLength       int64   `json:"audit_chain_length"`
	UnresolvedRPABarriers  int64   `json:"unresolved_rpa_barriers"`
	DatabaseSizeBytes      int64   `json:"database_size_bytes"`
	AverageCalcTimeMs      float64 `json:"average_calc_time_ms"`
}

type TelemetryCollector struct {
	db *sql.DB
}

func NewTelemetryCollector(db *sql.DB) *TelemetryCollector {
	return &TelemetryCollector{db: db}
}

func (tc *TelemetryCollector) CollectNodeTelemetry() (*SystemTelemetry, error) {
	st := &SystemTelemetry{}
	if tc.db == nil {
		return st, nil
	}
	_ = tc.db.QueryRow(`SELECT COUNT(*) FROM compliance_expedientes WHERE is_closed = 0;`).Scan(&st.TotalExpedientesActive)
	return st, nil
}
