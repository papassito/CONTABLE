package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidLinesCount = errors.New("FCOS_ERR_JOURNAL: Un asiento contable requiere al menos dos líneas para ser procesado")
	ErrNegativeAmount    = errors.New("FCOS_ERR_JOURNAL: Queda prohibido inyectar montos negativos en el libro diario")
	ErrDoubleSidedLine   = errors.New("FCOS_ERR_JOURNAL: Una línea debe afectar únicamente el debe o el haber, nunca ambos simultáneamente")
	ErrZeroMove          = errors.New("FCOS_ERR_JOURNAL: Una línea contable no puede registrar montos nulos en ambos lados")
)

type EntryStatus string

const (
	StatusDraft  EntryStatus = "BORRADOR"
	StatusPosted EntryStatus = "CONTABILIZADO"
)

// JournalLine describe el detalle transaccional de un asiento contable.
type JournalLine struct {
	AccountID   string `json:"account_id"`
	Description string `json:"description"`
	Debit       Cents  `json:"debit_cents"`
	Credit      Cents  `json:"credit_cents"`
}

// JournalEntry es la representación lógica del asiento de diario.
type JournalEntry struct {
	ID           string        `json:"id"`
	TenantID     string        `json:"tenant_id"`
	Number       string        `json:"number"`
	Date         string        `json:"date"` // Formato YYYY-MM-DD
	Concept      string        `json:"concept"`
	Status       EntryStatus   `json:"status"`
	Lines        []JournalLine `json:"lines"`
	CreatedAt    time.Time     `json:"created_at"`
	ReversalOfID string        `json:"reversal_of_entry_id,omitempty"`
}

// Validate analiza rigurosamente las invariantes financieras obligatorias del asiento contable.
func (je *JournalEntry) Validate() error {
	if len(je.Lines) < 2 {
		return ErrInvalidLinesCount
	}

	var totalDebits Cents = 0
	var totalCredits Cents = 0
	var err error

	for _, line := range je.Lines {
		// 1. Prohibir montos negativos
		if line.Debit < 0 || line.Credit < 0 {
			return ErrNegativeAmount
		}

		// 2. Prohibir dobles afectaciones simultáneas en una única línea
		if line.Debit > 0 && line.Credit > 0 {
			return ErrDoubleSidedLine
		}

		// 3. Prohibir registros nulos (debe y haber en cero)
		if line.Debit == 0 && line.Credit == 0 {
			return ErrZeroMove
		}

		// Acumulaciones seguras contra desbordamiento int64
		totalDebits, err = totalDebits.Add(line.Debit)
		if err != nil {
			return err
		}
		totalCredits, err = totalCredits.Add(line.Credit)
		if err != nil {
			return err
		}
	}

	if totalDebits != totalCredits {
		return ErrUnbalancedJournal
	}
	return nil
}
