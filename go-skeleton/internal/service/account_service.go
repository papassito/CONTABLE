package service

import (
	"context"
	"errors"
	"fmt"

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
	if len(account.Code) == 0 {
		return errors.New("el código contable no puede estar vacío")
	}

	existing, err := s.accountRepo.GetByCode(ctx, account.Code)
	if err == nil && existing != nil {
		return fmt.Errorf("la cuenta contable con código %s ya existe", account.Code)
	}

	firstChar := account.Code[0]
	var expectedType domain.AccountType
	switch firstChar {
	case '1':
		expectedType = domain.AccountTypeActivo
	case '2':
		expectedType = domain.AccountTypePasivo
	case '3':
		expectedType = domain.AccountTypePatrimonio
	case '4':
		expectedType = domain.AccountTypeIngreso
	case '5':
		expectedType = domain.AccountTypeGasto
	case '6':
		expectedType = domain.AccountTypeCosto
	default:
		return errors.New("código contable inválido: debe iniciar con un dígito de 1 a 6")
	}

	if account.Type != expectedType {
		return fmt.Errorf("tipo de cuenta %s no coincide con el código %s", account.Type, account.Code)
	}

	if account.ParentID != nil && *account.ParentID != "" {
		parent, err := s.accountRepo.GetByID(ctx, *account.ParentID)
		if err != nil || parent == nil {
			return errors.New("cuenta padre no encontrada")
		}
		if parent.AcceptsMove {
			if parent.CurrentBal != 0 {
				return fmt.Errorf("no se puede convertir la cuenta padre %s en mayorizadora porque posee un saldo de %d centavos", parent.Code, parent.CurrentBal)
			}
			parent.AcceptsMove = false
			if err := s.accountRepo.Update(ctx, parent); err != nil {
				return err
			}
		}
	}

	account.Status = domain.AccountStatusActiva
	account.CurrentBal = 0
	return s.accountRepo.Create(ctx, account)
}

func (s *accountService) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	return s.accountRepo.GetByID(ctx, id)
}

func (s *accountService) GetAccountByCode(ctx context.Context, code string) (*domain.Account, error) {
	return s.accountRepo.GetByCode(ctx, code)
}

func (s *accountService) ListAccounts(ctx context.Context) ([]*domain.Account, error) {
	return s.accountRepo.List(ctx, nil)
}

func (s *accountService) UpdateAccount(ctx context.Context, account *domain.Account) error {
	return s.accountRepo.Update(ctx, account)
}

func (s *accountService) DisableAccount(ctx context.Context, id string) error {
	account, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if account == nil {
		return domain.ErrAccountNotFound
	}

	if account.CurrentBal != 0 {
		return errors.New("no se puede desactivar una cuenta con saldo diferente de cero")
	}

	account.Status = domain.AccountStatusInactiva
	return s.accountRepo.Update(ctx, account)
}
