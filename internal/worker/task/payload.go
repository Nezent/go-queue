package task

type SendVerificationEmailPayload struct {
	Email string
	Token string
}

type EmailPayload struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

type WebSocketPayload struct {
	JobID   string `json:"job_id"`
	JobType string `json:"job_type"`
	Status  string `json:"status"`
}
