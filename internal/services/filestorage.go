package services

import (
	. "github.com/StampWallet/backend/internal/database"
	"os"
)

type FileStorageService interface {
	CreateStub(user *User) (FileMetadata, error)
	GetData(id string) (*os.File, error)
	Upload(fileMetadata FileMetadata, data os.File, mimetype string) (string, error)
	Remove(fileMetadata FileMetadata) error
}

type FileStorageServiceImpl struct {
	basePath     string
	baseServices BaseServices
}
