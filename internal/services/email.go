package services

import (
	"fmt"
	"log"

	mail "github.com/wneessen/go-mail"

	. "github.com/StampWallet/backend/internal/config"
	. "github.com/StampWallet/backend/internal/utils"
)

type EmailService interface {
	Send(email string, subject string, body string) error
}

type EmailServiceImpl struct {
	mailClient mail.Client
	smtpConfig SMTPConfig
	logger     log.Logger
}

func (service *EmailServiceImpl) Send(email string, subject string, body string) error {
	msg := mail.NewMsg(
		mail.WithEncoding(mail.EncodingQP),
		mail.WithCharset(mail.CharsetUTF8),
	)
	msg.Subject(subject)
	err := msg.AddTo(email)
	if err != nil {
		return fmt.Errorf("%s failed to add recipient: %+v", CallerFilename(), err)
	}
	msg.SetBodyString(mail.TypeTextHTML, subject)
	err = service.mailClient.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("%s failed to send email: %+v", CallerFilename(), err)
	}
	return nil
}
