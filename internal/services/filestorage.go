package services

import (
	"errors"
	"io"
	"os"

	. "github.com/StampWallet/backend/internal/database"
)

var (
	ErrNoSuchFile      = errors.New("No such file")
	ErrFileNotUploaded = errors.New("File not uploaded")
	ErrInvalidMimeType = errors.New("Invalid Mimetype")
)

var AllowedMimeTypes = []string{
	"image/jpeg",
	"image/png",
	"image/gif",
	"image/webp",
}

type FileStorageService interface {
	CreateStub(user *User) (*FileMetadata, error)
	// generally not very useful, will be served by a static server either way
	// FileMetadata, mimetype, error
	GetData(id string) (*os.File, string, error)
	// TODO how to recive an os.File from gin? data perhaps should be changed to reader
	Upload(fileMetadata FileMetadata, data io.Reader, mimetype string) (*FileMetadata, error)
	Remove(fileMetadata FileMetadata) error
}

type FileStorageServiceImpl struct {
	basePath     string
	baseServices BaseServices
}

func CreateFileStorageServiceImpl(baseServices BaseServices, basePath string) (*FileStorageServiceImpl, error) {
	//TODO check if path is accessible etc
	return &FileStorageServiceImpl{
		basePath:     basePath,
		baseServices: baseServices,
	}, nil
}

func (service *FileStorageServiceImpl) CreateStub(user *User) (*FileMetadata, string, error) {
	return nil, "", nil
}

func (service *FileStorageServiceImpl) GetData(id string) (*os.File, error) {
	return nil, nil
}

// limit upload to ~1mb
func (service *FileStorageServiceImpl) Upload(fileMetadata FileMetadata, data io.Reader, mimetype string) (*FileMetadata, error) {
	return nil, nil
}

func (service *FileStorageServiceImpl) Remove(fileMetadata FileMetadata) error {
	return nil
}
