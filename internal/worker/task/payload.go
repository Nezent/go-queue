package task

import "time"

type SendVerificationEmailPayload struct {
	Email string
	Token string
}

type EmailPayload struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

type JobPayload struct {
	Priority string       `json:"priority"`
	RunAt    time.Time    `json:"run_at"`
	Attempts int          `json:"attempts"`
	JobType  string       `json:"job_type"`
	Status   string       `json:"status"`
	Payload  EmailPayload `json:"payload"`
}

type WebSocketPayload struct {
	JobID   string `json:"job_id"`
	JobType string `json:"job_type"`
	Status  string `json:"status"`
}
