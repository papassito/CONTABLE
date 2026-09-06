package domain

import "time"

type EntryStatus string

const (
	EntryStatusBorrador      EntryStatus = "BORRADOR"
	EntryStatusContabilizado EntryStatus = "CONTABILIZADO"
	EntryStatusAnulado       EntryStatus = "ANULADO"
)

type JournalEntry struct {
	ID          string        `json:"id"`
	Number      string        `json:"number"`
	Date        time.Time     `json:"date"`
	Concept     string        `json:"concept"`
	Reference   string        `json:"reference"`
	Status      EntryStatus   `json:"status"`
	Lines       []JournalLine `json:"lines"`
	TotalDebit  float64       `json:"total_debit"`
	TotalCredit float64       `json:"total_credit"`
	CreatedBy   string        `json:"created_by"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type JournalLine struct {
	ID             string    `json:"id"`
	JournalEntryID string    `json:"journal_entry_id"`
	AccountID      string    `json:"account_id"`
	AccountCode    string    `json:"account_code"`
	Description    string    `json:"description"`
	Debit          float64   `json:"debit"`
	Credit         float64   `json:"credit"`
	ThirdPartyID   *string   `json:"third_party_id"`
	CreatedAt      time.Time `json:"created_at"`
}
