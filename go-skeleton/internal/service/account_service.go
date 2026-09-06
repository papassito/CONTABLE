package service

import (
	"context"
	"github.com/klik/contable-fix/internal/domain"
	"github.com/klik/contable-fix/internal/repository"
)

type AccountService interface {
	CreateAccount(ctx context.Context, account *domain.Account) error
	GetAccount(ctx context.Context, id string) (*domain.Account, error)
	GetAccountByCode(ctx context.Context, code string) (*domain.Account, error)
	ListAccounts(ctx context.Context) ([]*domain.Account, error)
	UpdateAccount(ctx context.Context, account *domain.Account) error
	DisableAccount(ctx context.Context, id string) error
}

type accountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{
		accountRepo: repo,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, account *domain.Account) error {
	// TODO: Validar jerarquía, unicidad de código y guardar
	return nil
}

func (s *accountService) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	// TODO: Implementar búsqueda por ID
	return nil, nil
}

func (s *accountService) GetAccountByCode(ctx context.Context, code string) (*domain.Account, error) {
	// TODO: Implementar búsqueda por código contable
	return nil, nil
}

func (s *accountService) ListAccounts(ctx context.Context) ([]*domain.Account, error) {
	// TODO: Implementar listado
	return nil, nil
}

func (s *accountService) UpdateAccount(ctx context.Context, account *domain.Account) error {
	// TODO: Implementar actualización
	return nil
}

func (s *accountService) DisableAccount(ctx context.Context, id string) error {
	// TODO: Implementar desactivación
	return nil
}
