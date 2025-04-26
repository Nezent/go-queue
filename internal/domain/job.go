package domain

import (
	"github.com/google/uuid"
)

type Job struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"`
	Payload   any       `json:"payload"`
	Status    string    `json:"status"`
	Priority  string    `json:"priority"`
	Attempts  int       `json:"attempts"`
	RunAt     string    `json:"run_at"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}
