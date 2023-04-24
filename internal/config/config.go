package config

type SMTPConfig struct {
    ServerHostname string
    ServerPort uint16
    Username string
    Password string
    SenderEmail string 
}

type Config struct {
    DatabaseUrl string
    SmtpConfig SMTPConfig
    ServerPort uint16
    StoragePath string
    BackendDomain string
    EmailVerificationFrontendURL string
}
