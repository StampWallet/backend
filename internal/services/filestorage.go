package services

import (
    "os"
    . "github.com/StampWallet/backend/internal/database"
)

type FileStorageService interface {
    CreateStub(user User) (FileMetadata, error)
    GetData(id string) (*os.File, error)
    Upload(fileMetadata FileMetadata, data os.File, mimetype string) (string, error)
    Remove(fileMetadata FileMetadata) error
}

type FileStorageServiceImpl struct {
    basePath string
    baseServices BaseServices
}
