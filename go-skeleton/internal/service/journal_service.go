package service

import (
	"context"
	"github.com/klik/contable-fix/internal/domain"
	"github.com/klik/contable-fix/internal/repository"
)

type JournalService interface {
	CreateDraft(ctx context.Context, entry *domain.JournalEntry) (*domain.JournalEntry, error)
	PostEntry(ctx context.Context, entryID string) error
	ReverseEntry(ctx context.Context, entryID string, reason string) (*domain.JournalEntry, error)
	GetEntry(ctx context.Context, id string) (*domain.JournalEntry, error)
	ListEntries(ctx context.Context, filter map[string]interface{}) ([]*domain.JournalEntry, error)
}

type journalService struct {
	journalRepo repository.JournalRepository
	accountRepo repository.AccountRepository
	ledgerRepo  repository.LedgerRepository
}

func NewJournalService(
	journalRepo repository.JournalRepository,
	accountRepo repository.AccountRepository,
	ledgerRepo repository.LedgerRepository,
) JournalService {
	return &journalService{
		journalRepo: journalRepo,
		accountRepo: accountRepo,
		ledgerRepo:  ledgerRepo,
	}
}

func (s *journalService) CreateDraft(ctx context.Context, entry *domain.JournalEntry) (*domain.JournalEntry, error) {
	// TODO: Validar partida doble (débito == crédito) y guardar borrador
	return nil, nil
}

func (s *journalService) PostEntry(ctx context.Context, entryID string) error {
	// TODO: Validar estado, asentar en libro mayor y marcar contabilizado
	return nil
}

func (s *journalService) ReverseEntry(ctx context.Context, entryID string, reason string) (*domain.JournalEntry, error) {
	// TODO: Generar contrapartida de reversión automática
	return nil, nil
}

func (s *journalService) GetEntry(ctx context.Context, id string) (*domain.JournalEntry, error) {
	// TODO: Consultar asiento con sus líneas
	return nil, nil
}

func (s *journalService) ListEntries(ctx context.Context, filter map[string]interface{}) ([]*domain.JournalEntry, error) {
	// TODO: Listar asientos filtrados
	return nil, nil
}
