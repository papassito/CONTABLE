package infrastructure

import (
	"database/sql"
	"fmt"
)

type SystemTelemetry struct {
	TotalExpedientesActive int64   `json:"total_expedientes_active"`
	TotalCalculationsCount int64   `json:"total_calculations_count"`
	AuditChainLength       int64   `json:"audit_chain_length"`
	UnresolvedRPABarriers  int64   `json:"unresolved_rpa_barriers"`
	DatabaseSizeBytes      int64   `json:"database_size_bytes"`
	AverageCalcTimeMs      float64 `json:"average_calc_time_ms"`
}

type TelemetryCollector interface {
	CollectNodeTelemetry() (*SystemTelemetry, error)
}

type telemetryCollector struct {
	db *sql.DB
}

func NewTelemetryCollector(db *sql.DB) TelemetryCollector {
	return &telemetryCollector{db: db}
}

func (tc *telemetryCollector) CollectNodeTelemetry() (*SystemTelemetry, error) {
	st := &SystemTelemetry{}

	err := tc.db.QueryRow(`SELECT COUNT(*) FROM compliance_expedientes WHERE is_closed = FALSE;`).Scan(&st.TotalExpedientesActive)
	if err != nil {
		return nil, fmt.Errorf("error en metrica expedientes: %w", err)
	}

	err = tc.db.QueryRow(`SELECT COUNT(*) FROM compliance_calculations;`).Scan(&st.TotalCalculationsCount)
	if err != nil {
		return nil, fmt.Errorf("error en metrica calculos: %w", err)
	}

	err = tc.db.QueryRow(`SELECT COALESCE(MAX(sequence_id), 0) FROM audit_hash_chain;`).Scan(&st.AuditChainLength)
	if err != nil {
		return nil, fmt.Errorf("error en metrica audit chain: %w", err)
	}

	err = tc.db.QueryRow(`SELECT COUNT(*) FROM compliance_expedientes WHERE current_stage = 'REQUIRES_HUMAN';`).Scan(&st.UnresolvedRPABarriers)
	if err != nil {
		return nil, fmt.Errorf("error en metrica rpa barriers: %w", err)
	}

	return st, nil
}
