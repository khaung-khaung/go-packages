package tests

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/banyar/go-packages/pkg/adapters"
	"github.com/banyar/go-packages/pkg/config"
	"github.com/banyar/go-packages/pkg/entities"
	"github.com/banyar/go-packages/pkg/frontlog"
	"go.uber.org/zap"
)

func TestLoadConfig(t *testing.T) {
	cfg := config.LoadConfig()
	if cfg == nil {
		t.Error("Expected configuration to be loaded, got nil")
	}

	// Initialize logger
	logConfig, logLevel := cfg.GetLoggingConfig()
	if err := frontlog.InitLogger(logConfig, logLevel); err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer frontlog.Logger.Sync()

	// // Now use the logger anywhere in your app
	frontlog.Logger.Info("Application started", zap.Any("", logConfig))

}

func TestSMTPSend(t *testing.T) {

	DSNSmtp := entities.DSNSmtp{
		SmtpServer:   os.Getenv("SMTP_SERVER"),
		SmtpPort:     os.Getenv("SMTP_PORT"),
		SmtpUser:     os.Getenv("SMTP_USER"),
		SmtpPassword: os.Getenv("SMTP_PASSWORD"),
	}

	tmpl := config.GetMailTemplete()
	data := entities.EmailData{
		Title: "Node Tag Process Summary",
		Items: []entities.Item{
			{Status: "success", Count: 100, FileName: "cloudshare.csv", Link: "https://www.google.com/"},
			{Status: "skipped", Count: 200, FileName: "cloudshare.csv", Link: "https://www.google.com/"},
			{Status: "failed", Count: 20, FileName: "cloudshare.csv", Link: "https://www.google.com/"},
		},
	}

	var body bytes.Buffer
	err := tmpl.Execute(&body, data)
	if err != nil {
		fmt.Printf("error executing template: %v", err)
	}

	mailUser := os.Getenv("MAIL_USERS")
	from := DSNSmtp.SmtpUser
	to := strings.Split(mailUser, ",")
	subject := "Subject: Node Tag Process Summary\n"
	contentType := "Content-Type: text/html; charset=UTF-8\n\n"
	msg := []byte(subject + contentType + body.String())
	smtpAdapter := adapters.NewSMTPMailAdapter(&DSNSmtp)
	smtpAdapter.MailService.SendMail(from, to, msg)
}
