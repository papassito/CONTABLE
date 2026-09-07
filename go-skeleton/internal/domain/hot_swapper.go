package domain

import (
	"context"
	"errors"
)

// HotSwapper maneja el reemplazo en caliente de reglas normativas sin recompilar
type HotSwapper interface {
	Swap(ctx context.Context, moduleName string) error
}

type hotSwapper struct{}

func NewHotSwapper() HotSwapper {
	return &hotSwapper{}
}

func (h *hotSwapper) Swap(ctx context.Context, moduleName string) error {
	if moduleName == "" {
		return errors.New("el nombre del módulo no puede estar vacío")
	}
	return nil
}
