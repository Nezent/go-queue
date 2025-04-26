package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"net/smtp"

	"github.com/Nezent/go-queue/internal/worker/task"

	"github.com/hibiken/asynq"
)

type TaskProcessor struct {
	// Inject any service or config (e.g., SMTP config)
	SMTPHost string
	SMTPPort string
	Auth     smtp.Auth
	From     string
}

func NewTaskProcessor(config SMTPConfig) *TaskProcessor {
	return &TaskProcessor{
		SMTPHost: config.Host,
		SMTPPort: config.Port,
		Auth:     smtp.PlainAuth("", config.Username, config.Password, config.Host),
		From:     config.From,
	}
}

func (p *TaskProcessor) HandleSendVerificationEmail(ctx context.Context, t *asynq.Task) error {
	var payload task.SendVerificationEmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	from := mail.Address{Name: "Sirajum Munir", Address: p.From}
	to := mail.Address{Address: payload.Email}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = "Verify your email"

	// Build the message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + fmt.Sprintf(
		"Click the link to verify your email:\nhttp://localhost:8080/api/v1/auth/verify?token=%s",
		payload.Token,
	)

	return smtp.SendMail(
		p.SMTPHost+":"+p.SMTPPort,
		p.Auth,
		p.From,               // still needs to be plain email address
		[]string{to.Address}, // list of recipient emails
		[]byte(message),
	)
}
