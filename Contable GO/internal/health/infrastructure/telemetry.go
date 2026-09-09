package infrastructure

import (
	"context"
	"database/sql"
	"time"
)

type NodeTelemetry struct {
	TotalExpedientesActive int64
	TotalCalculationsCount int64
	AuditChainLength       int64
	UnresolvedRPABarriers  int64
}

type TelemetryCollector struct {
	db          *sql.DB
	lastChecked time.Time
}

func NewTelemetryCollector(db *sql.DB) *TelemetryCollector {
	return &TelemetryCollector{
		db:          db,
		lastChecked: time.Now().UTC(),
	}
}

func (t *TelemetryCollector) CollectNodeTelemetry(ctx context.Context) (NodeTelemetry, error) {
	var telemetry NodeTelemetry
	t.lastChecked = time.Now().UTC()

	err := t.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM compliance_expedientes WHERE is_closed = 0").Scan(&telemetry.TotalExpedientesActive)
	if err != nil {
		return telemetry, err
	}

	err = t.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM compliance_calculations").Scan(&telemetry.TotalCalculationsCount)
	if err != nil {
		return telemetry, err
	}

	err = t.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM audit_hash_chain").Scan(&telemetry.AuditChainLength)
	if err != nil {
		return telemetry, err
	}

	err = t.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM compliance_expedientes WHERE current_stage = 'REQUIRES_HUMAN'").Scan(&telemetry.UnresolvedRPABarriers)
	if err != nil {
		return telemetry, err
	}

	return telemetry, nil
}

func (t *TelemetryCollector) CollectData(ctx context.Context) error {
	var err error
	t.lastChecked = time.Now().UTC()

	// Uso de '=' para evitar 'no new variables on left side of :='
	err = t.checkEngineHealth(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (t *TelemetryCollector) checkEngineHealth(_ context.Context) error {
	return nil
}
