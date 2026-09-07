package domain

import "time"

// AutomationTask representa la especificación y el estado de una tarea ejecutada por la infraestructura de RPA.
type AutomationTask struct {
	ID             string     `json:"id"`
	QueueName      string     `json:"queue_name"`
	PayloadJSON    string     `json:"payload_json"`
	Status         string     `json:"status"` // PENDING, RUNNING, COMPLETED, FAILED
	ExecutionError string     `json:"execution_error,omitempty"`
	Attempts       int        `json:"attempts"`
	MaxAttempts    int        `json:"max_attempts"`
	CreatedAtUTC   time.Time  `json:"created_at_utc"`
	StartedAtUTC   *time.Time `json:"started_at_utc,omitempty"`
	CompletedAtUTC *time.Time `json:"completed_at_utc,omitempty"`
}

// CanRetry determina si la tarea de automatización es elegible para reintento ante un fallo.
func (t *AutomationTask) CanRetry() bool {
	return t.Status == "FAILED" && t.Attempts < t.MaxAttempts
}
