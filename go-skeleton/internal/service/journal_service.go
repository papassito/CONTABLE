package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/gowebpki/jcs"
	"github.com/klik/fcos-kernel/internal/domain"
	"github.com/klik/fcos-kernel/internal/repository"
	"github.com/klik/fcos-kernel/pkg/validator"
)

// UnitOfWork define el puerto para coordinar transacciones ACID.
type UnitOfWork interface {
	Execute(ctx context.Context, fn func(ctx context.Context) error) error
}

type AuditService interface {
	AppendEvent(ctx context.Context, tenantID, eventType, actorID string, payload any) error
}

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
	uow         UnitOfWork
	auditSvc    AuditService
}

func NewJournalService(
	journalRepo repository.JournalRepository,
	accountRepo repository.AccountRepository,
	ledgerRepo repository.LedgerRepository,
	uow UnitOfWork,
	auditSvc AuditService,
) JournalService {
	return &journalService{
		journalRepo: journalRepo,
		accountRepo: accountRepo,
		ledgerRepo:  ledgerRepo,
		uow:         uow,
		auditSvc:    auditSvc,
	}
}

func (s *journalService) CreateDraft(ctx context.Context, entry *domain.JournalEntry) (*domain.JournalEntry, error) {
	if len(entry.Lines) < 2 {
		return nil, domain.ErrEmptyJournalLines
	}

	if !validator.ValidateDoubleEntry(entry.Lines) {
		return nil, domain.ErrUnbalancedJournal
	}

	// Validar que cada cuenta exista, esté activa y acepte movimientos
	for _, line := range entry.Lines {
		account, err := s.accountRepo.GetByID(ctx, line.AccountID)
		if err != nil || account == nil {
			return nil, domain.ErrAccountNotFound
		}
		if account.Status != domain.AccountStatusActiva {
			return nil, domain.ErrAccountInactive
		}
		if !account.AcceptsMove {
			return nil, domain.ErrAccountHasChildren
		}
	}

	entry.Status = domain.EntryStatusBorrador

	var totalDebit, totalCredit int64
	for _, line := range entry.Lines {
		totalDebit += line.Debit
		totalCredit += line.Credit
	}
	entry.TotalDebit = totalDebit
	entry.TotalCredit = totalCredit
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()

	err := s.journalRepo.Create(ctx, entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (s *journalService) PostEntry(ctx context.Context, entryID string) error {
	entry, err := s.journalRepo.GetByID(ctx, entryID)
	if err != nil {
		return err
	}
	if entry == nil {
		return domain.ErrJournalNotFound
	}

	if entry.Status != domain.EntryStatusBorrador {
		return domain.ErrEntryAlreadyPosted
	}

	lines, err := s.journalRepo.GetLinesByEntryID(ctx, entryID)
	if err != nil {
		return err
	}
	if len(lines) == 0 {
		lines = entry.Lines
	}
	if len(lines) < 2 {
		return domain.ErrEmptyJournalLines
	}

	if !validator.ValidateDoubleEntry(lines) {
		return domain.ErrUnbalancedJournal
	}

	return s.uow.Execute(ctx, func(txCtx context.Context) error {
		// Validar Periodo cerrado o declarado
		year, month, _ := entry.Date.Date()
		status, err := s.journalRepo.CheckPeriodStatus(txCtx, year, int(month))
		if err != nil {
			return err
		}
		if status == "CLOSED" || status == "DECLARED" {
			return domain.ErrClosedPeriod
		}

		ledgerEntries := make([]domain.LedgerEntry, len(lines))
		for i, line := range lines {
			ledgerEntries[i] = domain.LedgerEntry{
				AccountID:   line.AccountID,
				AccountCode: line.AccountCode,
				EntryDate:   entry.Date,
				Debit:       line.Debit,
				Credit:      line.Credit,
				Reference:   entry.Number,
			}
		}

		err = s.ledgerRepo.RecordMovements(txCtx, ledgerEntries)
		if err != nil {
			return err
		}

		for _, line := range lines {
			account, err := s.accountRepo.GetByID(txCtx, line.AccountID)
			if err != nil {
				return err
			}
			var amount int64
			switch account.Type {
			case domain.AccountTypeActivo, domain.AccountTypeGasto, domain.AccountTypeCosto:
				amount = line.Debit - line.Credit
			case domain.AccountTypePasivo, domain.AccountTypePatrimonio, domain.AccountTypeIngreso:
				amount = line.Credit - line.Debit
			default:
				amount = line.Debit - line.Credit
			}
			err = s.accountRepo.UpdateBalance(txCtx, line.AccountID, amount)
			if err != nil {
				return err
			}
		}

		err = s.journalRepo.UpdateStatus(txCtx, entryID, domain.EntryStatusContabilizado)
		if err != nil {
			return err
		}

		// Generar y emitir evento de auditoría canónico bajo RFC 8785
		auditPayload := map[string]any{
			"entry_id":    entry.ID,
			"number":      entry.Number,
			"total_debit": entry.TotalDebit,
			"status":      string(domain.EntryStatusContabilizado),
		}
		canonicalHash, err := computeCanonicalHash(auditPayload)
		if err != nil {
			return err
		}
		auditPayload["canonical_hash"] = canonicalHash

		if s.auditSvc != nil {
			err = s.auditSvc.AppendEvent(txCtx, "default-tenant", "POST_CONTABILIZAR", "system-user", auditPayload)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *journalService) ReverseEntry(ctx context.Context, entryID string, reason string) (*domain.JournalEntry, error) {
	var reversedEntry *domain.JournalEntry
	err := s.uow.Execute(ctx, func(txCtx context.Context) error {
		original, err := s.journalRepo.GetByID(txCtx, entryID)
		if err != nil {
			return err
		}
		if original == nil {
			return domain.ErrJournalNotFound
		}

		originalLines, err := s.journalRepo.GetLinesByEntryID(txCtx, entryID)
		if err != nil {
			return err
		}
		if len(originalLines) == 0 {
			originalLines = original.Lines
		}

		reversedLines := make([]domain.JournalLine, len(originalLines))
		for i, line := range originalLines {
			reversedLines[i] = domain.JournalLine{
				AccountID:    line.AccountID,
				AccountCode:  line.AccountCode,
				Description:  "Reversión: " + line.Description,
				Debit:        line.Credit,
				Credit:       line.Debit,
				ThirdPartyID: line.ThirdPartyID,
			}
		}

		reversedEntry = &domain.JournalEntry{
			Number:      "REV-" + original.Number,
			Date:        time.Now(),
			Concept:     "Reversión de " + original.Number + " - Motivo: " + reason,
			Reference:   original.Number,
			Status:      domain.EntryStatusBorrador,
			Lines:       reversedLines,
			TotalDebit:  original.TotalCredit,
			TotalCredit: original.TotalDebit,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err = s.journalRepo.Create(txCtx, reversedEntry)
		if err != nil {
			return err
		}

		// Generar y emitir evento de auditoría canónico bajo RFC 8785
		auditPayload := map[string]any{
			"original_entry_id": original.ID,
			"reversed_entry_id": reversedEntry.ID,
			"reason":            reason,
			"status":            string(domain.EntryStatusBorrador),
		}
		canonicalHash, err := computeCanonicalHash(auditPayload)
		if err != nil {
			return err
		}
		auditPayload["canonical_hash"] = canonicalHash

		if s.auditSvc != nil {
			err = s.auditSvc.AppendEvent(txCtx, "default-tenant", "VOID_ANULAR", "system-user", auditPayload)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return reversedEntry, nil
}

func (s *journalService) GetEntry(ctx context.Context, id string) (*domain.JournalEntry, error) {
	entry, err := s.journalRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	lines, err := s.journalRepo.GetLinesByEntryID(ctx, id)
	if err != nil {
		return nil, err
	}
	entry.Lines = lines
	return entry, nil
}

func (s *journalService) ListEntries(ctx context.Context, filter map[string]interface{}) ([]*domain.JournalEntry, error) {
	return s.journalRepo.List(ctx, filter)
}

func computeCanonicalHash(payload any) (string, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	canonical, err := jcs.Transform(payloadJSON)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(canonical)
	return hex.EncodeToString(hash[:]), nil
}
