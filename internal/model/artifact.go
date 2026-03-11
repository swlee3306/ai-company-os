package model

import "time"

type Artifact struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	Title     string         `json:"title"`
	ProjectID string         `json:"project_id,omitempty"`
	TaskID    string         `json:"task_id,omitempty"`
	URI       string         `json:"uri"`
	CreatedAt time.Time      `json:"created_at"`
	Meta      map[string]any `json:"meta,omitempty"`
}
