package services

import (
    "log"
    . "github.com/StampWallet/backend/internal/config"
    . "github.com/StampWallet/backend/internal/database"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type BaseServices struct {
    Logger *log.Logger
    Database GormDB
}

func getDatabase(config Config) (GormDB, error) {
    db, err := gorm.Open(postgres.Open(config.DatabaseUrl), &gorm.Config{})
    if db != nil && err != nil {
        return &GormDBImpl{Db: db}, nil
    } else {
        return nil, err
    }
}
