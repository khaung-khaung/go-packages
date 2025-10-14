package services

import (
	"github.com/banyar/go-packages/pkg/repositories"
)

type SMTPMailService struct {
	SmtpRepo *repositories.SMTPRepository
}

func NewSMTPService(smtpRepo *repositories.SMTPRepository) *SMTPMailService {
	return &SMTPMailService{
		SmtpRepo: smtpRepo,
	}
}

func (c *SMTPMailService) SendMail(from string, to []string, body []byte) {
	c.SmtpRepo.SendMail(from, to, body)
}
