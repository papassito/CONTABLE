package domain

import "time"

// WorkerJob define una tarea orquestada en segundo plano para el motor RPA
type WorkerJob struct {
	ID               string     `json:"id"`
	TenantID         string     `json:"tenant_id"`
	JobType          string     `json:"job_type"`
	Status           string     `json:"status"`
	RetryCount       int        `json:"retry_count"`
	MaxRetries       int        `json:"max_retries"`
	LastErrorMessage string     `json:"last_error_message"`
	Payload          string     `json:"payload"`
	ScheduledAtUTC   time.Time  `json:"scheduled_at_utc"`
	ExecutedAtUTC    *time.Time `json:"executed_at_utc"`
	CreatedAtUTC     time.Time  `json:"created_at_utc"`
}
