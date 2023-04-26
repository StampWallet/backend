package services

import (
	. "github.com/StampWallet/backend/internal/config"
	mail "github.com/wneessen/go-mail"
	"log"
)

type EmailService interface {
	Send(email string, subject string, body string) error
}

type EmailServiceImpl struct {
	mailClient mail.Client
	smtpConfig SMTPConfig
	logger     log.Logger
}
