package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/gowebpki/jcs"
	"github.com/klik/fcos-kernel/internal/domain"
	"github.com/klik/fcos-kernel/internal/service"
)

type mockJournalRepo struct {
	entries      map[string]*domain.JournalEntry
	lines        map[string][]domain.JournalLine
	periodStatus string
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
func (m *mockJournalRepo) CheckPeriodStatus(ctx context.Context, year int, month int) (string, error) {
	if m.periodStatus != "" {
		return m.periodStatus, nil
	}
	return "OPEN", nil
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
		a.CurrentBal += amount
		return nil
	}
	return errors.New("not found")
}

type mockLedgerRepo struct {
	entries []domain.LedgerEntry
	fail    bool
}

func (m *mockLedgerRepo) RecordMovements(ctx context.Context, entries []domain.LedgerEntry) error {
	if m.fail {
		return errors.New("ledger error")
	}
	m.entries = append(m.entries, entries...)
	return nil
}
func (m *mockLedgerRepo) GetMovementsByAccount(ctx context.Context, accountID string, from, to time.Time) ([]domain.LedgerEntry, error) {
	return nil, nil
}
func (m *mockLedgerRepo) GetTrialBalance(ctx context.Context, from, to time.Time) ([]domain.TrialBalanceItem, error) {
	return nil, nil
}

type mockUOW struct {
	wasRolledBack bool
}

func (m *mockUOW) Execute(ctx context.Context, fn func(txCtx context.Context) error) error {
	err := fn(ctx)
	if err != nil {
		m.wasRolledBack = true
	}
	return err
}

type mockAuditService struct {
	events []map[string]any
}

func (m *mockAuditService) AppendEvent(ctx context.Context, tenantID, eventType, actorID string, payload any) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = jcs.Transform(payloadJSON)
	if err != nil {
		return err
	}
	m.events = append(m.events, map[string]any{
		"tenantID":  tenantID,
		"eventType": eventType,
		"actorID":   actorID,
		"payload":   payload,
	})
	return nil
}

func TestCreateDraft_Balanced(t *testing.T) {
	jRepo := &mockJournalRepo{entries: make(map[string]*domain.JournalEntry), lines: make(map[string][]domain.JournalLine)}
	aRepo := &mockAccountRepo{accounts: make(map[string]*domain.Account)}
	lRepo := &mockLedgerRepo{}
	uow := &mockUOW{}
	audit := &mockAuditService{}

	// Configurar cuentas de prueba auxiliares activas
	aRepo.accounts["1"] = &domain.Account{ID: "1", Code: "110505", Status: domain.AccountStatusActiva, AcceptsMove: true}
	aRepo.accounts["2"] = &domain.Account{ID: "2", Code: "210505", Status: domain.AccountStatusActiva, AcceptsMove: true}

	svc := service.NewJournalService(jRepo, aRepo, lRepo, uow, audit)

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
	audit := &mockAuditService{}

	svc := service.NewJournalService(jRepo, aRepo, lRepo, uow, audit)

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
	audit := &mockAuditService{}

	entry := &domain.JournalEntry{
		ID:     "entry-1",
		Number: "001",
		Status: domain.EntryStatusContabilizado,
	}
	jRepo.entries["entry-1"] = entry

	svc := service.NewJournalService(jRepo, aRepo, lRepo, uow, audit)

	err := svc.PostEntry(context.Background(), "entry-1")
	if !errors.Is(err, domain.ErrEntryAlreadyPosted) {
		t.Errorf("Se esperaba ErrEntryAlreadyPosted, obtenido %v", err)
	}
}

func TestPostEntry_ClosedPeriod(t *testing.T) {
	jRepo := &mockJournalRepo{entries: make(map[string]*domain.JournalEntry), lines: make(map[string][]domain.JournalLine)}
	aRepo := &mockAccountRepo{accounts: make(map[string]*domain.Account)}
	lRepo := &mockLedgerRepo{}
	uow := &mockUOW{}
	audit := &mockAuditService{}

	entry := &domain.JournalEntry{
		ID:     "entry-1",
		Number: "001",
		Date:   time.Now(),
		Status: domain.EntryStatusBorrador,
		Lines: []domain.JournalLine{
			{AccountID: "1", Debit: 1000, Credit: 0},
			{AccountID: "2", Debit: 0, Credit: 1000},
		},
	}
	jRepo.entries["entry-1"] = entry
	jRepo.periodStatus = "CLOSED"

	svc := service.NewJournalService(jRepo, aRepo, lRepo, uow, audit)

	err := svc.PostEntry(context.Background(), "entry-1")
	if !errors.Is(err, domain.ErrClosedPeriod) {
		t.Errorf("Se esperaba ErrClosedPeriod, obtenido %v", err)
	}
}

func TestPostEntry_RollbackOnLedgerError(t *testing.T) {
	jRepo := &mockJournalRepo{entries: make(map[string]*domain.JournalEntry), lines: make(map[string][]domain.JournalLine)}
	aRepo := &mockAccountRepo{accounts: make(map[string]*domain.Account)}
	lRepo := &mockLedgerRepo{fail: true}
	uow := &mockUOW{}
	audit := &mockAuditService{}

	entry := &domain.JournalEntry{
		ID:     "entry-1",
		Number: "001",
		Date:   time.Now(),
		Status: domain.EntryStatusBorrador,
		Lines: []domain.JournalLine{
			{AccountID: "1", Debit: 1000, Credit: 0},
			{AccountID: "2", Debit: 0, Credit: 1000},
		},
	}
	jRepo.entries["entry-1"] = entry

	// Configurar cuentas auxiliares activas
	aRepo.accounts["1"] = &domain.Account{ID: "1", Code: "110505", Status: domain.AccountStatusActiva, AcceptsMove: true}
	aRepo.accounts["2"] = &domain.Account{ID: "2", Code: "210505", Status: domain.AccountStatusActiva, AcceptsMove: true}

	svc := service.NewJournalService(jRepo, aRepo, lRepo, uow, audit)

	err := svc.PostEntry(context.Background(), "entry-1")
	if err == nil {
		t.Fatal("Se esperaba error al fallar RecordMovements")
	}
	if !uow.wasRolledBack {
		t.Error("Se esperaba que la transacción hiciera rollback")
	}
}

func TestPostEntry_RollbackOnBalanceUpdateError(t *testing.T) {
	jRepo := &mockJournalRepo{entries: make(map[string]*domain.JournalEntry), lines: make(map[string][]domain.JournalLine)}
	aRepo := &mockAccountRepo{accounts: make(map[string]*domain.Account)}
	lRepo := &mockLedgerRepo{}
	uow := &mockUOW{}
	audit := &mockAuditService{}

	entry := &domain.JournalEntry{
		ID:     "entry-1",
		Number: "001",
		Date:   time.Now(),
		Status: domain.EntryStatusBorrador,
		Lines: []domain.JournalLine{
			{AccountID: "1", Debit: 1000, Credit: 0},
			{AccountID: "2", Debit: 0, Credit: 1000},
		},
	}
	jRepo.entries["entry-1"] = entry

	// Cuentas de prueba no configuradas en aRepo -> provocará error de "not found" en UpdateBalance
	// provocando que falle la transacción.

	svc := service.NewJournalService(jRepo, aRepo, lRepo, uow, audit)

	err := svc.PostEntry(context.Background(), "entry-1")
	if err == nil {
		t.Fatal("Se esperaba error de balance al no encontrar cuentas")
	}
	if !uow.wasRolledBack {
		t.Error("Se esperaba que la transacción hiciera rollback por error de balance")
	}
}

func TestPostEntry_Success(t *testing.T) {
	jRepo := &mockJournalRepo{entries: make(map[string]*domain.JournalEntry), lines: make(map[string][]domain.JournalLine)}
	aRepo := &mockAccountRepo{accounts: make(map[string]*domain.Account)}
	lRepo := &mockLedgerRepo{}
	uow := &mockUOW{}
	audit := &mockAuditService{}

	entry := &domain.JournalEntry{
		ID:     "entry-1",
		Number: "001",
		Date:   time.Now(),
		Status: domain.EntryStatusBorrador,
		Lines: []domain.JournalLine{
			{AccountID: "1", AccountCode: "110505", Debit: 1000, Credit: 0},
			{AccountID: "2", AccountCode: "210505", Debit: 0, Credit: 1000},
		},
	}
	jRepo.entries["entry-1"] = entry

	// Configurar cuentas auxiliares activas
	aRepo.accounts["1"] = &domain.Account{ID: "1", Code: "110505", Type: domain.AccountTypeActivo, Status: domain.AccountStatusActiva, AcceptsMove: true}
	aRepo.accounts["2"] = &domain.Account{ID: "2", Code: "210505", Type: domain.AccountTypePasivo, Status: domain.AccountStatusActiva, AcceptsMove: true}

	svc := service.NewJournalService(jRepo, aRepo, lRepo, uow, audit)

	err := svc.PostEntry(context.Background(), "entry-1")
	if err != nil {
		t.Fatalf("Error inesperado en PostEntry: %v", err)
	}

	if entry.Status != domain.EntryStatusContabilizado {
		t.Errorf("Se esperaba estado CONTABILIZADO, obtenido %s", entry.Status)
	}

	if aRepo.accounts["1"].CurrentBal != 1000 {
		t.Errorf("Se esperaba balance 1000 para cuenta 1, obtenido %d", aRepo.accounts["1"].CurrentBal)
	}
	if aRepo.accounts["2"].CurrentBal != 1000 {
		t.Errorf("Se esperaba balance 1000 para cuenta 2, obtenido %d", aRepo.accounts["2"].CurrentBal)
	}

	if len(lRepo.entries) != 2 {
		t.Errorf("Se esperaban 2 entradas de libro mayor, obtenido %d", len(lRepo.entries))
	}

	if len(audit.events) != 1 {
		t.Errorf("Se esperaba 1 evento de auditoría, obtenido %d", len(audit.events))
	}
}

func TestReverseEntry_Success(t *testing.T) {
	jRepo := &mockJournalRepo{entries: make(map[string]*domain.JournalEntry), lines: make(map[string][]domain.JournalLine)}
	aRepo := &mockAccountRepo{accounts: make(map[string]*domain.Account)}
	lRepo := &mockLedgerRepo{}
	uow := &mockUOW{}
	audit := &mockAuditService{}

	entry := &domain.JournalEntry{
		ID:     "entry-1",
		Number: "001",
		Status: domain.EntryStatusContabilizado,
		Lines: []domain.JournalLine{
			{AccountID: "1", AccountCode: "110505", Debit: 1000, Credit: 0, Description: "Debito original"},
			{AccountID: "2", AccountCode: "210505", Debit: 0, Credit: 1000, Description: "Credito original"},
		},
	}
	jRepo.entries["entry-1"] = entry

	svc := service.NewJournalService(jRepo, aRepo, lRepo, uow, audit)

	rev, err := svc.ReverseEntry(context.Background(), "entry-1", "Error de digitacion")
	if err != nil {
		t.Fatalf("Error inesperado en ReverseEntry: %v", err)
	}

	if rev.Number != "REV-001" {
		t.Errorf("Se esperaba número REV-001, obtenido %s", rev.Number)
	}

	if rev.Status != domain.EntryStatusBorrador {
		t.Errorf("Se esperaba estado BORRADOR para el asiento de reversión, obtenido %s", rev.Status)
	}

	if len(rev.Lines) != 2 {
		t.Fatalf("Se esperaban 2 líneas en reversión, obtenido %d", len(rev.Lines))
	}

	if rev.Lines[0].Debit != 0 || rev.Lines[0].Credit != 1000 {
		t.Errorf("Primera línea mal revertida: Debit %d, Credit %d", rev.Lines[0].Debit, rev.Lines[0].Credit)
	}
	if rev.Lines[1].Debit != 1000 || rev.Lines[1].Credit != 0 {
		t.Errorf("Segunda línea mal revertida: Debit %d, Credit %d", rev.Lines[1].Debit, rev.Lines[1].Credit)
	}
}
