package domain

import (
	"context"
	"time"
)

type IdempotencyStatus string

const (
	IdempotencyStatusProcessing IdempotencyStatus = "PROCESSING"
	IdempotencyStatusCompleted  IdempotencyStatus = "COMPLETED"
	IdempotencyStatusFailed     IdempotencyStatus = "FAILED"
)

// IdempotentRequest previene la doble ejecución de comandos mutativos
type IdempotentRequest struct {
	ID           string            `json:"id"`
	Status       IdempotencyStatus `json:"status"`
	ResponseJSON string            `json:"response_json,omitempty"`
	ProcessedAt  time.Time         `json:"processed_at"`
	ExpiresAt    time.Time         `json:"expires_at"`
}

// IsExpired determina si la solicitud idempotente ha expirado según el tiempo actual UTC.
func (r *IdempotentRequest) IsExpired() bool {
	return time.Now().UTC().After(r.ExpiresAt)
}

// CanExecute determina si la solicitud puede proceder a ejecutarse.
// Si ya expiró, o si falló previamente, se permite un reintento limpio.
// Si está actualmente en proceso o ya se completó con éxito, no se debe volver a ejecutar.
func (r *IdempotentRequest) CanExecute() bool {
	return r.IsExpired() || r.Status == IdempotencyStatusFailed
}

// MarkAsProcessing cambia el estado a PROCESSING actualizando la marca de tiempo en UTC.
func (r *IdempotentRequest) MarkAsProcessing() {
	r.Status = IdempotencyStatusProcessing
	r.ProcessedAt = time.Now().UTC()
}

// MarkAsCompleted finaliza la solicitud con éxito persistiendo la respuesta JSON.
func (r *IdempotentRequest) MarkAsCompleted(responseJSON string) {
	r.Status = IdempotencyStatusCompleted
	r.ResponseJSON = responseJSON
	r.ProcessedAt = time.Now().UTC()
}

// MarkAsFailed marca la ejecución como fallida permitiendo reintentos posteriores.
func (r *IdempotentRequest) MarkAsFailed() {
	r.Status = IdempotencyStatusFailed
	r.ProcessedAt = time.Now().UTC()
}

// IdempotencyRepository define el contrato para persistir y recuperar solicitudes idempotentes.
type IdempotencyRepository interface {
	FindByID(ctx context.Context, id string) (*IdempotentRequest, error)
	Save(ctx context.Context, req *IdempotentRequest) error
}
