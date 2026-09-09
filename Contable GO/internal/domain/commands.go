package domain

import "context"

// =============================================================================
// COMMAND CONTRACTS (Escritura / Mutación)
// =============================================================================

type CreateExpedienteCommand struct {
	TenantID        string `json:"tenant_id" validate:"required,uuid"`
	TaxpayerID      string `json:"taxpayer_id" validate:"required,uuid"`
	PeriodYear      int    `json:"period_year" validate:"required,gte=2000"`
	PeriodMonth     int    `json:"period_month" validate:"required,min=1,max=12"`
	PeriodType      string `json:"period_type" validate:"required,oneof=MONTHLY BIMONTHLY ANNUAL CUSTOM"`
	AuthorityDomain string `json:"authority_domain" validate:"required,oneof=SAT IMSS ISRTP_ESTATAL"`
	ActorID         string `json:"actor_id" validate:"required,uuid"`
}

type ExecuteCalculationCommand struct {
	TenantID        string `json:"tenant_id" validate:"required,uuid"`
	ExpedienteID    string `json:"expediente_id" validate:"required,uuid"`
	RuleCode        string `json:"rule_code" validate:"required"`
	BaseAmountCents int64  `json:"base_amount_cents" validate:"required,gte=0"`
	PaymentDate     string `json:"payment_date" validate:"required,datetime=2006-01-02"`
	ActorID         string `json:"actor_id" validate:"required,uuid"`
}

type TransitionWorkflowCommand struct {
	TenantID     string `json:"tenant_id" validate:"required,uuid"`
	ExpedienteID string `json:"expediente_id" validate:"required,uuid"`
	TargetStage  string `json:"target_stage" validate:"required"`
	Reason       string `json:"reason" validate:"required"`
	ActorID      string `json:"actor_id" validate:"required,uuid"`
}

// Interfaces de Handlers de Comandos
type CommandHandler[C any, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

// =============================================================================
// QUERY CONTRACTS (Lectura / Búsquedas)
// =============================================================================

type GetExpedienteByIDQuery struct {
	TenantID     string `json:"tenant_id"`
	ExpedienteID string `json:"expediente_id"`
}

// Interfaces de Handlers de Consultas
type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, query Q) (R, error)
}

// =============================================================================
// RESPONSE CONTRACTS
// =============================================================================
type ExpedienteResponse struct {
	ExpedienteID   string `json:"expediente_id"`
	ExpedienteCode string `json:"expediente_code"`
	CurrentStage   string `json:"current_stage"`
}
