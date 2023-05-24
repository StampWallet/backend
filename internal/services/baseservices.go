package services

import (
	"log"

	. "github.com/StampWallet/backend/internal/config"
	. "github.com/StampWallet/backend/internal/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Struct with most commonly used services
type BaseServices struct {
	Logger   *log.Logger
	Database GormDB
}

// Opens database connection
func GetDatabase(config Config) (GormDB, error) {
	db, err := gorm.Open(postgres.Open(config.DatabaseUrl), &gorm.Config{})
	if db != nil && err == nil {
		return &GormDBImpl{Db: db}, nil
	} else {
		return nil, err
	}
}

// Creates a new logger with prefix
func NewPrefix(logger *log.Logger, prefix string) *log.Logger {
	if prefix[len(prefix)-1] != ' ' {
		prefix += " "
	}
	newLogger := log.New(logger.Writer(), prefix, 0)
	newLogger.SetFlags(logger.Flags())
	return newLogger
}

// Creates a new logger with prefix from BaseServices
func (b BaseServices) NewPrefix(prefix string) BaseServices {
	return BaseServices{
		Logger:   NewPrefix(b.Logger, prefix),
		Database: b.Database,
	}
}
