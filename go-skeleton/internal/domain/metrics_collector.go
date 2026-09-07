package domain

import (
	"sync"
	"time"
)

// MetricsCollector define la interfaz para el registro de eventos y rendimiento del nodo local.
type MetricsCollector interface {
	IncCounter(metricName string, value int64)
	RecordLatency(metricName string, duration time.Duration)
	GetCounter(metricName string) int64
	GetAverageLatency(metricName string) time.Duration
}

type latencyStats struct {
	totalSum   time.Duration
	totalCount int64
}

type metricsCollector struct {
	mu        sync.RWMutex
	counters  map[string]int64
	latencies map[string]*latencyStats
}

// NewMetricsCollector inicializa una implementación segura para hilos del colector de métricas.
func NewMetricsCollector() MetricsCollector {
	return &metricsCollector{
		counters:  make(map[string]int64),
		latencies: make(map[string]*latencyStats),
	}
}

func (m *metricsCollector) IncCounter(metricName string, value int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[metricName] += value
}

func (m *metricsCollector) RecordLatency(metricName string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	stats, exists := m.latencies[metricName]
	if !exists {
		stats = &latencyStats{}
		m.latencies[metricName] = stats
	}
	stats.totalSum += duration
	stats.totalCount++
}

func (m *metricsCollector) GetCounter(metricName string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.counters[metricName]
}

func (m *metricsCollector) GetAverageLatency(metricName string) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	stats, exists := m.latencies[metricName]
	if !exists || stats.totalCount == 0 {
		return 0
	}
	return stats.totalSum / time.Duration(stats.totalCount)
}
