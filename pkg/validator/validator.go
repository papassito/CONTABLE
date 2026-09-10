package validator

import (
	"github.com/klik/fcos-kernel/internal/domain"
)

func ValidateDoubleEntry(lines []domain.JournalLine) bool {
	var totalDebit domain.Cents
	var totalCredit domain.Cents
	var err error

	for _, line := range lines {
		totalDebit, err = totalDebit.Add(line.Debit)
		if err != nil {
			return false
		}
		totalCredit, err = totalCredit.Add(line.Credit)
		if err != nil {
			return false
		}
	}
	return totalDebit == totalCredit
}
