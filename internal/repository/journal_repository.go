package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/klik/fcos-kernel/internal/domain"
)

type JournalRepository interface {
	Create(ctx context.Context, entry *domain.JournalEntry) error
	GetByID(ctx context.Context, id string) (*domain.JournalEntry, error)
	GetByNumber(ctx context.Context, number string) (*domain.JournalEntry, error)
	List(ctx context.Context, filter map[string]interface{}) ([]*domain.JournalEntry, error)
	UpdateStatus(ctx context.Context, id string, status domain.EntryStatus) error
	GetLinesByEntryID(ctx context.Context, entryID string) ([]domain.JournalLine, error)
	CheckPeriodStatus(ctx context.Context, year int, month int) (string, error)
}

type journalEntryRow struct {
	ID          string
	Number      string
	Date        string
	Concept     string
	Reference   string
	Status      string
	TotalDebit  int64
	TotalCredit int64
}

type journalLineRow struct {
	ID             string
	JournalEntryID string
	AccountID      string
	AccountCode    string
	Description    string
	Debit          int64
	Credit         int64
	ThirdPartyID   sql.NullString
}

func toDomainEntry(row *journalEntryRow, lines []domain.JournalLine) *domain.JournalEntry {
	if row == nil {
		return nil
	}
	return &domain.JournalEntry{
		ID:        row.ID,
		TenantID:  "", // Se integrará dinámicamente en la Fase 2
		Number:    row.Number,
		Date:      row.Date,
		Concept:   row.Concept,
		Status:    domain.EntryStatus(row.Status),
		Lines:     lines,
		CreatedAt: time.Now(), // Fallback seguro para compatibilidad
	}
}

func fromDomainEntry(entry *domain.JournalEntry) (*journalEntryRow, []journalLineRow) {
	if entry == nil {
		return nil, nil
	}

	var totalDebit, totalCredit int64
	var lineRows []journalLineRow

	for _, line := range entry.Lines {
		totalDebit += int64(line.Debit)
		totalCredit += int64(line.Credit)

		lineRows = append(lineRows, journalLineRow{
			ID:             uuid.New().String(),
			JournalEntryID: entry.ID,
			AccountID:      line.AccountID,
			AccountCode:    "", // Se resolverá en la Fase 2 mediante consulta si es necesario
			Description:    line.Description,
			Debit:          int64(line.Debit),
			Credit:         int64(line.Credit),
			ThirdPartyID:   sql.NullString{},
		})
	}

	return &journalEntryRow{
		ID:          entry.ID,
		Number:      entry.Number,
		Date:        entry.Date,
		Concept:     entry.Concept,
		Reference:   "",
		Status:      string(entry.Status),
		TotalDebit:  totalDebit,
		TotalCredit: totalCredit,
	}, lineRows
}

type sqlJournalRepository struct {
	db *sql.DB
}

func NewJournalRepository(db *sql.DB) JournalRepository {
	return &sqlJournalRepository{db: db}
}

func (r *sqlJournalRepository) Create(ctx context.Context, entry *domain.JournalEntry) error {
	exec := getExecutor(ctx, r.db)
	row, lineRows := fromDomainEntry(entry)
	_, err := exec.ExecContext(ctx, "INSERT INTO journal_entries (id, number, date, concept, reference, status, total_debit, total_credit) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		row.ID, row.Number, row.Date, row.Concept, row.Reference, row.Status, row.TotalDebit, row.TotalCredit)
	if err != nil {
		return err
	}
	for _, line := range lineRows {
		_, err = exec.ExecContext(ctx, "INSERT INTO journal_lines (id, journal_entry_id, account_id, account_code, description, debit, credit, third_party_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			line.ID, line.JournalEntryID, line.AccountID, line.AccountCode, line.Description, line.Debit, line.Credit, line.ThirdPartyID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *sqlJournalRepository) GetByID(ctx context.Context, id string) (*domain.JournalEntry, error) {
	exec := getExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, "SELECT id, number, date, concept, reference, status, total_debit, total_credit FROM journal_entries WHERE id = ?", id)
	var dbRow journalEntryRow
	err := row.Scan(&dbRow.ID, &dbRow.Number, &dbRow.Date, &dbRow.Concept, &dbRow.Reference, &dbRow.Status, &dbRow.TotalDebit, &dbRow.TotalCredit)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	lines, err := r.GetLinesByEntryID(ctx, dbRow.ID)
	if err != nil {
		return nil, err
	}
	return toDomainEntry(&dbRow, lines), nil
}

func (r *sqlJournalRepository) GetByNumber(ctx context.Context, number string) (*domain.JournalEntry, error) {
	exec := getExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, "SELECT id, number, date, concept, reference, status, total_debit, total_credit FROM journal_entries WHERE number = ?", number)
	var dbRow journalEntryRow
	err := row.Scan(&dbRow.ID, &dbRow.Number, &dbRow.Date, &dbRow.Concept, &dbRow.Reference, &dbRow.Status, &dbRow.TotalDebit, &dbRow.TotalCredit)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	lines, err := r.GetLinesByEntryID(ctx, dbRow.ID)
	if err != nil {
		return nil, err
	}
	return toDomainEntry(&dbRow, lines), nil
}

func (r *sqlJournalRepository) List(ctx context.Context, filter map[string]interface{}) ([]*domain.JournalEntry, error) {
	exec := getExecutor(ctx, r.db)
	rows, err := exec.QueryContext(ctx, "SELECT id, number, date, concept, reference, status, total_debit, total_credit FROM journal_entries")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.JournalEntry
	for rows.Next() {
		var dbRow journalEntryRow
		if err := rows.Scan(&dbRow.ID, &dbRow.Number, &dbRow.Date, &dbRow.Concept, &dbRow.Reference, &dbRow.Status, &dbRow.TotalDebit, &dbRow.TotalCredit); err != nil {
			return nil, err
		}
		lines, err := r.GetLinesByEntryID(ctx, dbRow.ID)
		if err != nil {
			return nil, err
		}
		list = append(list, toDomainEntry(&dbRow, lines))
	}
	return list, nil
}

func (r *sqlJournalRepository) UpdateStatus(ctx context.Context, id string, status domain.EntryStatus) error {
	exec := getExecutor(ctx, r.db)
	_, err := exec.ExecContext(ctx, "UPDATE journal_entries SET status = ? WHERE id = ?", status, id)
	return err
}

func (r *sqlJournalRepository) GetLinesByEntryID(ctx context.Context, entryID string) ([]domain.JournalLine, error) {
	exec := getExecutor(ctx, r.db)
	rows, err := exec.QueryContext(ctx, "SELECT id, journal_entry_id, account_id, account_code, description, debit, credit, third_party_id FROM journal_lines WHERE journal_entry_id = ?", entryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var lines []domain.JournalLine
	for rows.Next() {
		var dbRow journalLineRow
		if err := rows.Scan(&dbRow.ID, &dbRow.JournalEntryID, &dbRow.AccountID, &dbRow.AccountCode, &dbRow.Description, &dbRow.Debit, &dbRow.Credit, &dbRow.ThirdPartyID); err != nil {
			return nil, err
		}
		lines = append(lines, domain.JournalLine{
			AccountID:   dbRow.AccountID,
			Description: dbRow.Description,
			Debit:       domain.Cents(dbRow.Debit),
			Credit:      domain.Cents(dbRow.Credit),
		})
	}
	return lines, nil
}

func (r *sqlJournalRepository) CheckPeriodStatus(ctx context.Context, year int, month int) (string, error) {
	exec := getExecutor(ctx, r.db)
	var status string
	err := exec.QueryRowContext(ctx, "SELECT status FROM compliance_periods WHERE period_year = ? AND period_month = ?", year, month).Scan(&status)
	if err == sql.ErrNoRows {
		return "OPEN", nil
	}
	return status, err
}
