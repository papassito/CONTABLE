package domain

import "time"

// MetricsCollector define la interfaz para el registro de eventos y rendimiento del nodo local.
type MetricsCollector interface {
	IncCounter(metricName string, value int64)
	RecordLatency(metricName string, duration time.Duration)
}

type metricsCollector struct{}

func NewMetricsCollector() MetricsCollector {
	return &metricsCollector{}
}

func (m *metricsCollector) IncCounter(metricName string, value int64) {
	// Base implementation
}

func (m *metricsCollector) RecordLatency(metricName string, duration time.Duration) {
	// Base implementation
}
