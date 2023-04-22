package services

import (
    "log"
    mail "github.com/wneessen/go-mail"
)

type SMTPConfig struct {
    ServerHostname string
    ServerPort uint16
    Username string
    Password string
    SenderEmail string 
}

type EmailService interface {
    Send(email string, subject string, body string) error
}

type EmailServiceImpl struct {
    mailClient mail.Client
    smtpConfig SMTPConfig
    logger log.Logger
}
