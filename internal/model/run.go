package model

type Run struct {
	ID         string `json:"id"`
	TaskID     string `json:"task_id"`
	RunnerType string `json:"runner_type"`
	Pipeline   string `json:"pipeline"`
	Status     string `json:"status"` // running|done|failed
	StartedAt  string `json:"started_at"`
	EndedAt    string `json:"ended_at,omitempty"`
	Summary    string `json:"summary,omitempty"`
}
