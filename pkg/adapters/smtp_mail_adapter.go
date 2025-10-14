package adapters

import (
	entities "github.com/banyar/go-packages/pkg/entities"
	"github.com/banyar/go-packages/pkg/interfaces"
	"github.com/banyar/go-packages/pkg/repositories"
	"github.com/banyar/go-packages/pkg/services"
)

type SMTPMailAdapter struct {
	MailService interfaces.IMailService
}

func NewSMTPMailAdapter(DSNSmtp *entities.DSNSmtp) *SMTPMailAdapter {
	smtpRepo := repositories.ConnectSMTP(DSNSmtp)
	return &SMTPMailAdapter{
		MailService: services.NewSMTPService(smtpRepo),
	}
}
