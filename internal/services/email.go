package services

import (
	"log"

	. "github.com/StampWallet/backend/internal/config"
	mail "github.com/wneessen/go-mail"
)

type EmailService interface {
	Send(email string, subject string, body string) error
}

type EmailServiceImpl struct {
	mailClient mail.Client
	smtpConfig SMTPConfig
	logger     *log.Logger
}

func (service *EmailServiceImpl) Send(email string, subject string, body string) error {
	return nil
}

func CreateEmailServiceImpl(smtpConfig SMTPConfig, logger *log.Logger) (*EmailServiceImpl, error) {
	client, err := mail.NewClient(
		smtpConfig.ServerHostname,
		//mail.WithLogger(*logger),
		mail.WithPassword(smtpConfig.Password),
		mail.WithUsername(smtpConfig.Username),
		mail.WithPort(int(smtpConfig.ServerPort)),
		mail.WithSSL())
	if err != nil {
		return nil, err
	}
	return &EmailServiceImpl{
		mailClient: *client,
		smtpConfig: smtpConfig,
		logger:     logger,
	}, nil
}
