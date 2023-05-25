package services

import (
	"fmt"
	"log"

	mail "github.com/wneessen/go-mail"

	. "github.com/StampWallet/backend/internal/config"
	. "github.com/StampWallet/backend/internal/utils"
)

// An EmailService is a service for sending emails.
// It's only a thin wrapper over github.com/wneessen/go-mail
// Configuration options are documented in config.SMTPConfig.
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

		// authorization
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithPassword(smtpConfig.Password),
		mail.WithUsername(smtpConfig.Username),

		// connection options - ssl is currently forced
		mail.WithPort(int(smtpConfig.ServerPort)),
		mail.WithSSL(),

		// not 100% sure if this is necessary
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

// Sends an email message to address email, with subject and body from args.
// body is currently assumed to be valid html - watch out for injections.
// It's recommended to use html templates to create emails,
// but this struct is not responsible for this.
func (service *EmailServiceImpl) Send(email string, subject string, body string) error {
	msg := mail.NewMsg(
		// URL encode - should make this email readable on text only clients
		// (but that was not tested yet)
		mail.WithEncoding(mail.EncodingQP),
		// UTF-8
		mail.WithCharset(mail.CharsetUTF8),
	)

	// set up message content
	msg.Subject(subject)
	msg.From(service.smtpConfig.SenderEmail)
	err := msg.AddTo(email)
	if err != nil {
		return fmt.Errorf("%s failed to add recipient: %+v", CallerFilename(), err)
	}
	msg.SetBodyString(mail.TypeTextHTML, body)

	err = service.mailClient.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("%s failed to send email: %+v", CallerFilename(), err)
	}
	return nil
}
