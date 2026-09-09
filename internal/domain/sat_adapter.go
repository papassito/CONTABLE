package domain

import (
	"context"
)

type SATAdapter interface {
	Download32DOpinion(ctx context.Context, rfc string, secretPayload []byte) (*SAT32DResponseDTO, error)
}
