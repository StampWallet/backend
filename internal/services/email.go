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
	logger     *log.Logger
}

func CreateEmailServiceImpl(smtpConfig SMTPConfig, logger *log.Logger) (*EmailServiceImpl, error) {
	client, err := mail.NewClient(
		smtpConfig.ServerHostname,
		//mail.WithLogger(*logger),
		mail.WithPassword(smtpConfig.Password),
		mail.WithUsername(smtpConfig.Username),
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithPort(int(smtpConfig.ServerPort)),
		mail.WithSSL(),
		mail.WithoutNoop())
	if err != nil {
		return nil, err
	}
	return &EmailServiceImpl{
		mailClient: *client,
		smtpConfig: smtpConfig,
		logger:     logger,
	}, nil
}

func (service *EmailServiceImpl) Send(email string, subject string, body string) error {
	msg := mail.NewMsg(
		mail.WithEncoding(mail.EncodingQP),
		mail.WithCharset(mail.CharsetUTF8),
	)
	msg.Subject(subject)
	msg.From(service.smtpConfig.SenderEmail)
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
