package validator

import (
	"math"
	"github.com/klik/contable-fix/internal/domain"
)

func ValidateDoubleEntry(lines []domain.JournalLine) bool {
	var totalDebit, totalCredit float64
	for _, line := range lines {
		totalDebit += line.Debit
		totalCredit += line.Credit
	}
	return math.Abs(totalDebit-totalCredit) < 0.0001
}
