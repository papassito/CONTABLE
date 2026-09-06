package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/klik/contable-fix/internal/domain"
	"github.com/klik/contable-fix/internal/service"
)

type mockJournalRepo struct {
	entries map[string]*domain.JournalEntry
	lines   map[string][]domain.JournalLine
}

func (m *mockJournalRepo) Create(ctx context.Context, entry *domain.JournalEntry) error {
	m.entries[entry.ID] = entry
	return nil
}
func (m *mockJournalRepo) GetByID(ctx context.Context, id string) (*domain.JournalEntry, error) {
	return m.entries[id], nil
}
func (m *mockJournalRepo) GetByNumber(ctx context.Context, number string) (*domain.JournalEntry, error) {
	for _, e := range m.entries {
		if e.Number == number {
			return e, nil
		}
	}
	return nil, nil
}
func (m *mockJournalRepo) List(ctx context.Context, filter map[string]interface{}) ([]*domain.JournalEntry, error) {
	var list []*domain.JournalEntry
	for _, e := range m.entries {
		list = append(list, e)
	}
	return list, nil
}
func (m *mockJournalRepo) UpdateStatus(ctx context.Context, id string, status domain.EntryStatus) error {
	if e, ok := m.entries[id]; ok {
		e.Status = status
		return nil
	}
	return errors.New("not found")
}
func (m *mockJournalRepo) GetLinesByEntryID(ctx context.Context, entryID string) ([]domain.JournalLine, error) {
	return m.lines[entryID], nil
}

type mockAccountRepo struct {
	accounts map[string]*domain.Account
}

func (m *mockAccountRepo) Create(ctx context.Context, account *domain.Account) error {
	m.accounts[account.ID] = account
	return nil
}
func (m *mockAccountRepo) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	return m.accounts[id], nil
}
func (m *mockAccountRepo) GetByCode(ctx context.Context, code string) (*domain.Account, error) {
	for _, a := range m.accounts {
		if a.Code == code {
			return a, nil
		}
	}
	return nil, nil
}
func (m *mockAccountRepo) List(ctx context.Context, filter map[string]interface{}) ([]*domain.Account, error) {
	var list []*domain.Account
	for _, a := range m.accounts {
		list = append(list, a)
	}
	return list, nil
}
func (m *mockAccountRepo) Update(ctx context.Context, account *domain.Account) error {
	m.accounts[account.ID] = account
	return nil
}
func (m *mockAccountRepo) Delete(ctx context.Context, id string) error {
	delete(m.accounts, id)
	return nil
}
func (m *mockAccountRepo) UpdateBalance(ctx context.Context, id string, amount int64) error {
	if a, ok := m.accounts[id]; ok {
		a.CurrentBal += float64(amount)
		return nil
	}
	return errors.New("not found")
}

type mockLedgerRepo struct {
	entries []domain.LedgerEntry
}

func (m *mockLedgerRepo) RecordMovements(ctx context.Context, entries []domain.LedgerEntry) error {
	m.entries = append(m.entries, entries...)
	return nil
}
func (m *mockLedgerRepo) GetMovementsByAccount(ctx context.Context, accountID string, from, to time.Time) ([]domain.LedgerEntry, error) {
	return nil, nil
}
func (m *mockLedgerRepo) GetTrialBalance(ctx context.Context, from, to time.Time) ([]domain.TrialBalanceItem, error) {
	return nil, nil
}

type mockUOW struct{}

func (m *mockUOW) Execute(ctx context.Context, fn func(txCtx context.Context) error) error {
	return fn(ctx)
}

func TestCreateDraft_Balanced(t *testing.T) {
	jRepo := &mockJournalRepo{entries: make(map[string]*domain.JournalEntry), lines: make(map[string][]domain.JournalLine)}
	aRepo := &mockAccountRepo{accounts: make(map[string]*domain.Account)}
	lRepo := &mockLedgerRepo{}
	uow := &mockUOW{}

	// Configurar cuentas de prueba auxiliares activas
	aRepo.accounts["1"] = &domain.Account{ID: "1", Code: "110505", Status: domain.AccountStatusActiva, AcceptsMove: true}
	aRepo.accounts["2"] = &domain.Account{ID: "2", Code: "210505", Status: domain.AccountStatusActiva, AcceptsMove: true}

	svc := service.NewJournalService(jRepo, aRepo, lRepo, uow)

	entry := &domain.JournalEntry{
		ID:     "entry-1",
		Number: "001",
		Lines: []domain.JournalLine{
			{AccountID: "1", Debit: 1000, Credit: 0},
			{AccountID: "2", Debit: 0, Credit: 1000},
		},
	}

	res, err := svc.CreateDraft(context.Background(), entry)
	if err != nil {
		t.Fatalf("Error inesperado en CreateDraft: %v", err)
	}
	if res.Status != domain.EntryStatusBorrador {
		t.Errorf("Se esperaba estado BORRADOR, obtenido %s", res.Status)
	}
}

func TestCreateDraft_Unbalanced(t *testing.T) {
	jRepo := &mockJournalRepo{entries: make(map[string]*domain.JournalEntry), lines: make(map[string][]domain.JournalLine)}
	aRepo := &mockAccountRepo{accounts: make(map[string]*domain.Account)}
	lRepo := &mockLedgerRepo{}
	uow := &mockUOW{}

	svc := service.NewJournalService(jRepo, aRepo, lRepo, uow)

	entry := &domain.JournalEntry{
		ID:     "entry-1",
		Number: "001",
		Lines: []domain.JournalLine{
			{AccountID: "1", Debit: 1000, Credit: 0},
			{AccountID: "2", Debit: 0, Credit: 999}, // Diferencia de un centavo
		},
	}

	_, err := svc.CreateDraft(context.Background(), entry)
	if !errors.Is(err, domain.ErrUnbalancedJournal) {
		t.Errorf("Se esperaba ErrUnbalancedJournal, obtenido %v", err)
	}
}

func TestPostEntry_AlreadyPosted(t *testing.T) {
	jRepo := &mockJournalRepo{entries: make(map[string]*domain.JournalEntry), lines: make(map[string][]domain.JournalLine)}
	aRepo := &mockAccountRepo{accounts: make(map[string]*domain.Account)}
	lRepo := &mockLedgerRepo{}
	uow := &mockUOW{}

	entry := &domain.JournalEntry{
		ID:     "entry-1",
		Number: "001",
		Status: domain.EntryStatusContabilizado,
	}
	jRepo.entries["entry-1"] = entry

	svc := service.NewJournalService(jRepo, aRepo, lRepo, uow)

	err := svc.PostEntry(context.Background(), "entry-1")
	if !errors.Is(err, domain.ErrEntryAlreadyPosted) {
		t.Errorf("Se esperaba ErrEntryAlreadyPosted, obtenido %v", err)
	}
}
