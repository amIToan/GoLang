package email

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error
}
type GmailSender struct {
	name              string
	FromEmailAddress  string
	FromEmailPassword string
}

func CreateNewEmailSender(name string, fromEmailAddress string, FromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		FromEmailAddress:  fromEmailAddress,
		FromEmailPassword: FromEmailPassword,
	}
}
func (gmailSender *GmailSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", gmailSender.name, gmailSender.FromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}
	smtpAuth := smtp.PlainAuth("", gmailSender.FromEmailAddress, gmailSender.FromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}
