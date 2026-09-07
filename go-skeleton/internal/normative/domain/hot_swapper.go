package domain

import "context"

// HotSwapper define el contrato específico dentro del módulo normativo para el intercambio de reglas en caliente.
type HotSwapper interface {
	SwapRule(ctx context.Context, ruleCode string, rawFormula string) error
	IsActive(ruleCode string) bool
}
