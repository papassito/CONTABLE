package domain

import (
	"context"
)

// NodeMetricsCollector define el comportamiento para recopilar y reportar la telemetría del hardware del nodo.
type NodeMetricsCollector interface {
	CollectCPUUsage(ctx context.Context) (float64, error)
	CollectMemoryUsage(ctx context.Context) (uint64, error)
}
