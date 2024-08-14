package entities

type DSNSmtp struct {
	SmtpServer   string //smtp server
	SmtpPort     string //smtp port
	SmtpUser     string //smtp user
	SmtpPassword string //smtp password
}

type Item struct {
	Status   string
	Count    int
	FileName string
	Link     string
}

// EmailData represents the data to be used in the email template
type EmailData struct {
	Title string
	Items []Item
}
