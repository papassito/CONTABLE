package domain

import "context"

// SATAutomationAdapter define el contrato del robot (RPA) para interactuar directamente con el portal del SAT.
type SATAutomationAdapter interface {
	ExecuteRPAOpinionDownload(ctx context.Context, rfc string, password string) ([]byte, error)
	DownloadCFDIs(ctx context.Context, rfc string, password string, year int, month int) ([][]byte, error)
}
