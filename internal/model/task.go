package model

import "time"

type Task struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	Desc             string    `json:"desc,omitempty"`
	State            string    `json:"state"`
	Assignee         string    `json:"assignee,omitempty"`
	ReviewerRequired bool      `json:"reviewer_required"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
