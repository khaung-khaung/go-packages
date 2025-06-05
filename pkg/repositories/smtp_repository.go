package repositories

import (
	"fmt"
	"net/smtp"

	entities "github.com/banyar/go-packages/pkg/entities"
)

type SMTPRepository struct {
	addr string
	auth smtp.Auth
}

func ConnectSMTP(DSNSmtp *entities.DSNSmtp) *SMTPRepository {
	return &SMTPRepository{
		addr: fmt.Sprintf("%s:%s", DSNSmtp.SmtpServer, DSNSmtp.SmtpPort),                     // SMTP server address
		auth: smtp.PlainAuth("", DSNSmtp.SmtpUser, DSNSmtp.SmtpPassword, DSNSmtp.SmtpServer), // Authentication
	}
}

func (r *SMTPRepository) SendMail(from string, to []string, body []byte) string {
	err := smtp.SendMail(
		r.addr,
		r.auth,
		from,
		to,
		body,
	)
	if err != nil {
		return "Failed to send email :" + err.Error()
	}
	return "Email sent successfully!"
}
