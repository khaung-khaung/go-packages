package interfaces

type IMailService interface {
	SendMail(from string, to []string, body []byte)
}
