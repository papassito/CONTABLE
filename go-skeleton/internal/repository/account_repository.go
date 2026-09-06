package repository

import (
	"context"
	"github.com/klik/contable-fix/internal/domain"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	GetByID(ctx context.Context, id string) (*domain.Account, error)
	GetByCode(ctx context.Context, code string) (*domain.Account, error)
	List(ctx context.Context, filter map[string]interface{}) ([]*domain.Account, error)
	Update(ctx context.Context, account *domain.Account) error
	Delete(ctx context.Context, id string) error
	UpdateBalance(ctx context.Context, id string, amount float64) error
}
