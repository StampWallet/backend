package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

// SMTP client config
type SMTPConfig struct {
	ServerHostname string // SMTP server hostname
	ServerPort     uint16 // SMTP server port
	Username       string // Authorization username
	Password       string // Authorization password
	SenderEmail    string // Email to use in the "From" field
}

type Config struct {
	DatabaseUrl                   string     // Database URL
	SmtpConfig                    SMTPConfig // SMTP Client config
	ListenIP                      string     // Hostname:port this server will listen on
	StoragePath                   string     // File storage path
	BackendURL                    string     // Public DNS domain this server is reachable from
	VerificationEmailSubject      string     // String with verification email subject
	VerificationEmailBodyTemplate string     // Template that receives the email verification token
	StaticPath                    string     // Static file path
}

// Returns config with default values
func GetDefaultConfig() Config {
	return Config{
		DatabaseUrl: "localhost",
		SmtpConfig: SMTPConfig{
			ServerHostname: "localhost",
			ServerPort:     465,
			Username:       "test",
			Password:       "test",
			SenderEmail:    "test@localhost",
		},
		ListenIP:                      "localhost:8080",
		StoragePath:                   "/tmp/",
		StaticPath:                    "static",
		BackendURL:                    "http://localhost:8080/",
		VerificationEmailSubject:      "email subject",
		VerificationEmailBodyTemplate: "http://localhost:8080/static/emailVerification.html?token={{ .Token}}",
	}
}

// Loads config from file under path
func LoadConfig(path string) (Config, error) {
	k := koanf.New(".")
	fileSuccess := true
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		fmt.Printf("failed to load config file: %+v\n", err)
		fmt.Printf("continuing\n")
		fileSuccess = false
		//return Config{}, err
	}

	prefix := "STAMPWALLET_"
	loaded := 0
	err := k.Load(env.Provider(prefix, ".", func(s string) string {
		loaded += 1
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, prefix)),
			"_", ".", -1,
		)
	}), nil)
	if err != nil {
		return Config{}, err
	}

	if loaded == 0 && !fileSuccess {
		return Config{}, errors.New("failed to load config both from env and file")
	}

	var config Config
	k.UnmarshalWithConf("", &config, koanf.UnmarshalConf{})
	return config, nil
}

// Saves config to file under path
func SaveConfig(cfg Config, path string) error {
	k := koanf.New(".")
	if err := k.Load(structs.Provider(cfg, "koanf"), nil); err != nil {
		return err
	}
	result, err := k.Marshal(yaml.Parser())
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	n, err := file.Write(result)
	defer file.Close()
	if n < len(result) {
		return errors.New("failed to write the whole config file. wrote only " + strconv.Itoa(n))
	} else if err != nil {
		return err
	}

	return nil
}
