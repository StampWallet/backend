package config

import (
    "github.com/StampWallet/backend/internal/services"
)

type Config struct {
    DatabaseUrl string
    SmtpConfig services.SMTPConfig
    ServerPort uint16
    StoragePath string
    BackendDomain string
    EmailVerificationFrontendURL string
}
