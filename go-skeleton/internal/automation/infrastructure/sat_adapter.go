package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/klik/fcos-kernel/internal/domain"
)

type SATAdapter interface {
	Download32DOpinion(ctx context.Context, rfc string, secretPayload []byte) (*domain.SAT32DResponseDTO, error)
}

type satAdapter struct {
	parser domain.DocumentParser
}

func NewSATAdapter(parser domain.DocumentParser) SATAdapter {
	return &satAdapter{parser: parser}
}

func (a *satAdapter) Download32DOpinion(ctx context.Context, rfc string, secretPayload []byte) (*domain.SAT32DResponseDTO, error) {
	if len(secretPayload) == 0 {
		return nil, fmt.Errorf("no existen credenciales en Secrets Vault para RFC %s", rfc)
	}
	docID := uuid.New().String()
	mockPDFContent := []byte("%PDF-1.4 ... OPINION DEL CUMPLIMIENTO: POSITIVA ...")
	return &domain.SAT32DResponseDTO{
		DocumentID:      docID,
		RFC:             rfc,
		OpinionResult:   "POSITIVA",
		PdfRawBytes:     mockPDFContent,
		PdfSHA256:       "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		DownloadedAtUTC: time.Now().UTC(),
	}, nil
}
