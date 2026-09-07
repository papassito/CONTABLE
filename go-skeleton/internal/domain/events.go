package domain

import (
	"time"
)

// BaseDomainEvent interfaz que todo evento del FCOS debe implementar
type BaseDomainEvent interface {
	EventID() string
	TenantID() string
	EventType() string
	OccurredAt() time.Time
}

// Evento: ExpedienteCreado
type ExpedienteCreadoEvent struct {
	ID              string    `json:"event_id"`
	TenantUUID      string    `json:"tenant_id"`
	ExpedienteUUID  string    `json:"expediente_id"`
	TaxpayerUUID    string    `json:"taxpayer_id"`
	AuthorityDomain string    `json:"authority_domain"`
	ExpedienteCode  string    `json:"expediente_code"`
	ActorUUID       string    `json:"actor_id"`
	TimestampUTC    time.Time `json:"timestamp_utc"`
}

func (e ExpedienteCreadoEvent) EventID() string { return e.ID }

func (e ExpedienteCreadoEvent) TenantID() string { return e.TenantUUID }

func (e ExpedienteCreadoEvent) EventType() string { return "ExpedienteCreado" }

func (e ExpedienteCreadoEvent) OccurredAt() time.Time { return e.TimestampUTC }

// Evento: CalculoEjecutado
type CalculoEjecutadoEvent struct {
	ID                    string    `json:"event_id"`
	TenantUUID            string    `json:"tenant_id"`
	ExpedienteUUID        string    `json:"expediente_id"`
	CalculationUUID       string    `json:"calculation_id"`
	RuleVersionUUID       string    `json:"rule_version_id"`
	TotalAmountCents      int64     `json:"total_amount_cents"`
	CalculationOutputHash string    `json:"calculation_output_hash"`
	ActorUUID             string    `json:"actor_id"`
	TimestampUTC          time.Time `json:"timestamp_utc"`
}

func (e CalculoEjecutadoEvent) EventID() string { return e.ID }

func (e CalculoEjecutadoEvent) TenantID() string { return e.TenantUUID }

func (e CalculoEjecutadoEvent) EventType() string { return "CalculoEjecutado" }

func (e CalculoEjecutadoEvent) OccurredAt() time.Time { return e.TimestampUTC }

// Evento: IntervencionHumanaRequerida
type IntervencionHumanaRequeridaEvent struct {
	ID             string    `json:"event_id"`
	TenantUUID     string    `json:"tenant_id"`
	ExpedienteUUID string    `json:"expediente_id"`
	WorkerID       string    `json:"worker_id"`
	PortalName     string    `json:"portal_name"`
	BarrierType    string    `json:"barrier_type"` // CAPTCHA, MFA, LAYOUT_CHANGE
	Message        string    `json:"message"`
	TimestampUTC   time.Time `json:"timestamp_utc"`
}

func (e IntervencionHumanaRequeridaEvent) EventID() string { return e.ID }

func (e IntervencionHumanaRequeridaEvent) TenantID() string { return e.TenantUUID }

func (e IntervencionHumanaRequeridaEvent) EventType() string { return "IntervencionHumanaRequerida" }

func (e IntervencionHumanaRequeridaEvent) OccurredAt() time.Time { return e.TimestampUTC }
