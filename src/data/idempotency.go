package domain

import "time"

// IdempotentRequest previene la doble ejecución de comandos mutativos
type IdempotentRequest struct {
	ID           string    `json:"id"`
	ResponseJSON string    `json:"response_json"`
	ProcessedAt  time.Time `json:"processed_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}
