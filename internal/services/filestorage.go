package services

import (
	"errors"
	"os"

	. "github.com/StampWallet/backend/internal/database"
)

var NoSuchFile = errors.New("No such file")
var FileNotUploaded = errors.New("File not uploaded")
var InvalidMimeType = errors.New("Invalid Mimetype")

var AllowedMimeTypes = []string{
	"image/jpeg",
	"image/png",
	"image/gif",
	"image/webp",
}

type FileStorageService interface {
	CreateStub(user *User) (*FileMetadata, error)
	// generally not very useful, will be served by a static server either way
	GetData(id string) (*os.File, error)
	// TODO how to recive an os.File from gin? data perhaps should be changed to reader
	Upload(fileMetadata FileMetadata, data *os.File, mimetype string) (*FileMetadata, error)
	Remove(fileMetadata FileMetadata) error
}

type FileStorageServiceImpl struct {
	basePath     string
	baseServices BaseServices
}

func (service *FileStorageServiceImpl) CreateStub(user *User) (*FileMetadata, error) {
	return nil, nil
}

func (service *FileStorageServiceImpl) GetData(id string) (*os.File, error) {
	return nil, nil
}

func (service *FileStorageServiceImpl) Upload(fileMetadata FileMetadata, data *os.File, mimetype string) (*FileMetadata, error) {
	return nil, nil
}

func (service *FileStorageServiceImpl) Remove(fileMetadata FileMetadata) error {
	return nil
}
