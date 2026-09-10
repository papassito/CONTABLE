package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gowebpki/jcs"
	"github.com/klik/fcos-kernel/internal/domain"
	"github.com/klik/fcos-kernel/internal/repository"
)

// UnitOfWork define el puerto para coordinar transacciones ACID.
type UnitOfWork interface {
	Execute(ctx context.Context, fn func(ctx context.Context) error) error
}

// AuditService define el puerto para persistir eventos firmados criptográficamente.
type AuditService interface {
	AppendEvent(ctx context.Context, tenantID, eventType, actorID string, payload any) error
}

// JournalService define los casos de uso transaccionales del Libro Diario.
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
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	if entry.TenantID == "" {
		return nil, ErrTenantRequired
	}

	if entry.TenantID != tenantID {
		return nil, ErrTenantMismatch
	}

	if err := entry.Validate(); err != nil {
		return nil, err
	}

	accountsCache := make(map[string]*domain.Account)
	for _, line := range entry.Lines {
		account, ok := accountsCache[line.AccountID]
		if !ok {
			var err error
			account, err = s.accountRepo.GetByID(ctx, tenantID, line.AccountID)
			if err != nil || account == nil {
				return nil, domain.ErrAccountNotFound
			}
			if account.TenantID != tenantID {
				return nil, ErrTenantMismatch
			}
			accountsCache[line.AccountID] = account
		}
		if account.Status != domain.AccountActive {
			return nil, domain.ErrAccountInactive
		}
		if !account.AcceptsMove {
			return nil, domain.ErrAccountHasChildren
		}
	}

	entry.Status = domain.StatusDraft
	entry.CreatedAt = time.Now()

	err = s.journalRepo.Create(ctx, entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (s *journalService) PostEntry(ctx context.Context, entryID string) error {
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return err
	}

	return s.uow.Execute(ctx, func(txCtx context.Context) error {
		entry, err := s.journalRepo.GetByID(txCtx, entryID)
		if err != nil {
			return err
		}
		if entry == nil {
			return domain.ErrJournalNotFound
		}
		if entry.TenantID != tenantID {
			return ErrTenantMismatch
		}
		if entry.Status != domain.StatusDraft {
			return domain.ErrEntryAlreadyPosted
		}

		lines, err := s.journalRepo.GetLinesByEntryID(txCtx, entryID)
		if err != nil {
			return err
		}
		if len(lines) == 0 {
			lines = entry.Lines
		}

		tempEntry := &domain.JournalEntry{Lines: lines}
		if err := tempEntry.Validate(); err != nil {
			return err
		}

		parts := strings.Split(entry.Date, "-")
		if len(parts) != 3 {
			return errors.New("FCOS_ERR_JOURNAL: formato de fecha contable inválido")
		}
		year, errY := strconv.Atoi(parts[0])
		month, errM := strconv.Atoi(parts[1])
		if errY != nil || errM != nil || month < 1 || month > 12 {
			return errors.New("FCOS_ERR_JOURNAL: formato de fecha contable inválido")
		}

		status, err := s.journalRepo.CheckPeriodStatus(txCtx, year, month)
		if err != nil {
			return err
		}
		if status == "CLOSED" || status == "DECLARED" {
			return domain.ErrClosedPeriod
		}

		entryDate, err := time.Parse("2006-01-02", entry.Date)
		if err != nil {
			return err
		}

		ledgerEntries := make([]domain.LedgerEntry, len(lines))
		netBalances := make(map[string]domain.Cents)
		accountsCache := make(map[string]*domain.Account)

		for i, line := range lines {
			account, ok := accountsCache[line.AccountID]
			if !ok {
				var err error
				account, err = s.accountRepo.GetByID(txCtx, tenantID, line.AccountID)
				if err != nil || account == nil {
					return domain.ErrAccountNotFound
				}
				if account.TenantID != tenantID {
					return ErrTenantMismatch
				}
				if account.Status != domain.AccountActive {
					return domain.ErrAccountInactive
				}
				if !account.AcceptsMove {
					return domain.ErrAccountHasChildren
				}
				accountsCache[line.AccountID] = account
			}

			ledgerEntries[i] = domain.LedgerEntry{
				AccountID:   line.AccountID,
				AccountCode: account.Code,
				EntryDate:   entryDate,
				Debit:       int64(line.Debit),
				Credit:      int64(line.Credit),
				Reference:   entry.Number,
			}

			delta, errSub := line.Debit.Sub(line.Credit)
			if errSub != nil {
				return errSub
			}
			var errAdd error
			netBalances[line.AccountID], errAdd = netBalances[line.AccountID].Add(delta)
			if errAdd != nil {
				return errAdd
			}
		}

		if err = s.ledgerRepo.RecordMovements(txCtx, ledgerEntries); err != nil {
			return err
		}

		for accountID, netAmount := range netBalances {
			if netAmount != 0 { // Solo actualiza si hubo impacto neto real
				if err = s.accountRepo.UpdateBalance(txCtx, tenantID, accountID, netAmount); err != nil {
					return err
				}
			}
		}

		if err = s.journalRepo.UpdateStatus(txCtx, entryID, domain.StatusPosted); err != nil {
			return err
		}

		if s.auditSvc != nil {
			var totalDebit int64
			for _, line := range lines {
				totalDebit += int64(line.Debit)
			}
			auditPayload := map[string]any{
				"entry_id":    entry.ID,
				"number":      entry.Number,
				"total_debit": totalDebit,
				"posted_at":   time.Now().UTC().Format(time.RFC3339),
				"status":      string(domain.StatusPosted),
			}
			canonicalHash, err := computeCanonicalHash(auditPayload)
			if err != nil {
				return err
			}
			auditPayload["canonical_hash"] = canonicalHash

			if err = s.auditSvc.AppendEvent(txCtx, tenantID, "POST_CONTABILIZAR", "system-user", auditPayload); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *journalService) ReverseEntry(ctx context.Context, entryID string, reason string) (*domain.JournalEntry, error) {
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var reversedEntry *domain.JournalEntry
	err = s.uow.Execute(ctx, func(txCtx context.Context) error {
		original, err := s.journalRepo.GetByID(txCtx, entryID)
		if err != nil {
			return err
		}
		if original == nil {
			return domain.ErrJournalNotFound
		}
		if original.TenantID != tenantID {
			return ErrTenantMismatch
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
				AccountID:   line.AccountID,
				Description: "Reversión: " + line.Description,
				Debit:       line.Credit,
				Credit:      line.Debit,
			}
		}

		reversedEntry = &domain.JournalEntry{
			ID:           uuid.New().String(),
			TenantID:     tenantID,
			Number:       "REV-" + original.Number,
			Date:         time.Now().Format("2006-01-02"),
			Concept:      "Reversión de " + original.Number + " - Motivo: " + reason,
			Status:       domain.StatusDraft,
			Lines:        reversedLines,
			CreatedAt:    time.Now(),
			ReversalOfID: original.ID,
		}

		err = s.journalRepo.Create(txCtx, reversedEntry)
		if err != nil {
			return err
		}

		var totalDebit int64
		for _, line := range reversedLines {
			totalDebit += int64(line.Debit)
		}

		auditPayload := map[string]any{
			"original_entry_id": original.ID,
			"reversed_entry_id": reversedEntry.ID,
			"reason":            reason,
			"reversed_at":       time.Now().UTC().Format(time.RFC3339),
			"status":            string(domain.StatusDraft),
		}
		canonicalHash, err := computeCanonicalHash(auditPayload)
		if err != nil {
			return err
		}
		auditPayload["canonical_hash"] = canonicalHash

		if s.auditSvc != nil {
			err = s.auditSvc.AppendEvent(txCtx, tenantID, "VOID_ANULAR", "system-user", auditPayload)
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
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	entry, err := s.journalRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	if entry.TenantID != tenantID {
		return nil, ErrTenantMismatch
	}
	lines, err := s.journalRepo.GetLinesByEntryID(ctx, id)
	if err != nil {
		return nil, err
	}
	entry.Lines = lines
	return entry, nil
}

func (s *journalService) ListEntries(ctx context.Context, filter map[string]interface{}) ([]*domain.JournalEntry, error) {
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	entries, err := s.journalRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	var filtered []*domain.JournalEntry
	for _, e := range entries {
		if e.TenantID == tenantID {
			filtered = append(filtered, e)
		}
	}
	return filtered, nil
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
