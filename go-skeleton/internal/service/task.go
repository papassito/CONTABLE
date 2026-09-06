package domain

import "time"

// WorkerTask modela un sub-paso atómico de ejecución dentro de un Job
type WorkerTask struct {
	ID         string     `json:"id"`
	JobID      string     `json:"job_id"`
	StepName   string     `json:"step_name"`
	Status     string     `json:"status"`
	OutputData string     `json:"output_data"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
}
