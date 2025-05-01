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
