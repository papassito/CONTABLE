package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/klik/fcos-kernel/internal/domain"
)

type SATAdapter interface {
	Download32DOpinion(ctx context.Context, rfc string, encryptedSecretPayload []byte) (*domain.SAT32DResponseDTO, error)
}

type satAdapter struct {
	parser DocumentParser
}

func NewSATAdapter(parser DocumentParser) SATAdapter {
	return &satAdapter{parser: parser}
}

func (a *satAdapter) Download32DOpinion(ctx context.Context, rfc string, secretPayload []byte) (*domain.SAT32DResponseDTO, error) {
	if len(secretPayload) == 0 {
		return nil, fmt.Errorf("no existen credenciales o llaves de e.firma en Secrets Vault para RFC %s", rfc)
	}

	// Simulación de interacción con el adaptador HTTP/TLS
	// En producción, el worker utiliza las cookies de sesión obtenidas con la e.firma

	// Caso A: Detección de Barrera CAPTCHA / MFA
	if ctx.Value("simulate_barrier") == true {
		return nil, &domain.RPABarrierError{
			WorkerID:    "worker-sat-01",
			PortalName:  "Portal SAT - Opinión de Cumplimiento 32-D",
			BarrierType: "CAPTCHA",
			CaptchaB64:  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
			Message:     "Se requiere resolución de CAPTCHA visual para acceder a la descarga",
		}
	}

	// Caso B: Descarga Exitosa
	docID := uuid.New().String()
	mockPDFContent := []byte("%PDF-1.4 ... OPINION DEL CUMPLIMIENTO DE OBLIGACIONES FISCALES: POSITIVA ...")

	return &domain.SAT32DResponseDTO{
		DocumentID:      docID,
		RFC:             rfc,
		OpinionResult:   "POSITIVA",
		PdfRawBytes:     mockPDFContent,
		PdfSHA256:       "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		DownloadedAtUTC: time.Now().UTC(),
	}, nil
}
