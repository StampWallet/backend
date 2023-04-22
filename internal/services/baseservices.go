package services

import (
    "log"
    "gorm.io/gorm"
)

type BaseServices struct {
    Logger *log.Logger
    Database *gorm.DB
}
