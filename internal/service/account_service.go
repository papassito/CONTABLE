package service

import (
	"context"
	"errors"

	"github.com/klik/fcos-kernel/internal/domain"
	"github.com/klik/fcos-kernel/internal/repository"
)

type contextKey string

const tenantIDKey contextKey = "tenant_id"

var (
	ErrTenantRequired = errors.New("FCOS_ERR_SECURITY: tenant_id is required")
	ErrTenantMismatch = errors.New("FCOS_ERR_SECURITY: tenant_id mismatch between context and request")
)

func getTenantID(ctx context.Context) (string, error) {
	tenantID, ok := ctx.Value(tenantIDKey).(string)
	if !ok || tenantID == "" {
		return "", ErrTenantRequired
	}
	return tenantID, nil
}

// WithTenantID asocia de forma segura la identidad del tenant al contexto utilizando la clave privada tipada.
func WithTenantID(ctx context.Context, tenantID string) (context.Context, error) {
	if tenantID == "" {
		return nil, ErrTenantRequired
	}
	return context.WithValue(ctx, tenantIDKey, tenantID), nil
}

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
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return err
	}

	if account.TenantID == "" {
		return ErrTenantRequired
	}

	if account.TenantID != tenantID {
		return ErrTenantMismatch
	}

	if len(account.Code) == 0 {
		return errors.New("el código contable no puede estar vacío")
	}

	existing, err := s.accountRepo.GetByCode(ctx, tenantID, account.Code)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("FCOS_ERR_ACCOUNT: account code already exists for this tenant")
	}

	if account.ParentID != "" {
		parent, err := s.accountRepo.GetByID(ctx, tenantID, account.ParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return errors.New("cuenta padre no encontrada")
		}

		accountsList, err := s.accountRepo.List(ctx, tenantID)
		if err != nil {
			return err
		}

		accountsMap := make(map[string]*domain.Account)
		for _, a := range accountsList {
			accountsMap[a.ID] = a
		}
		accountsMap[account.ID] = account

		if domain.DetectCycle(accountsMap, account.ID) {
			return domain.ErrHierarchyCycle
		}
	}

	account.Status = domain.AccountActive
	account.Balance = 0
	return s.accountRepo.Create(ctx, account)
}

func (s *accountService) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	return s.accountRepo.GetByID(ctx, tenantID, id)
}

func (s *accountService) GetAccountByCode(ctx context.Context, code string) (*domain.Account, error) {
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	return s.accountRepo.GetByCode(ctx, tenantID, code)
}

func (s *accountService) ListAccounts(ctx context.Context) ([]*domain.Account, error) {
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	return s.accountRepo.List(ctx, tenantID)
}

func (s *accountService) UpdateAccount(ctx context.Context, account *domain.Account) error {
	return errors.New("FCOS_ERR_NOT_IMPLEMENTED: UpdateAccount is not implemented in FCOS v2.2")
}

func (s *accountService) DisableAccount(ctx context.Context, id string) error {
	return errors.New("FCOS_ERR_NOT_IMPLEMENTED: DisableAccount is not implemented in FCOS v2.2")
}
