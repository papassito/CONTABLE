package domain

import (
	"context"

	sharedDomain "github.com/klik/fcos-kernel/internal/domain"
)

type SATAdapter interface {
	Download32DOpinion(ctx context.Context, rfc string, secretPayload []byte) (*sharedDomain.SAT32DResponseDTO, error)
}
