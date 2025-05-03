package domain

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID        uuid.UUID      `json:"id"`
	UserID    uuid.UUID      `json:"user_id"`
	Type      string         `json:"type"`
	Payload   map[string]any `json:"payload"`
	Status    string         `json:"status"`
	Priority  string         `json:"priority"`
	Attempts  int            `json:"attempts"`
	RunAt     time.Time      `json:"run_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type JobCreateRequestDTO struct {
	Type     string         `json:"type"`
	Payload  map[string]any `json:"payload"`
	Priority string         `json:"priority"`
	RunAt    string         `json:"run_at"`
}

type JobStatusResponseDTO struct {
	Type     string    `json:"type"`
	Status   string    `json:"status"`
	Priority string    `json:"priority"`
	Attempts int       `json:"attempts"`
	RunAt    time.Time `json:"run_at"`
}
