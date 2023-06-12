package services

import (
	"database/sql"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	. "github.com/StampWallet/backend/internal/database"
	"github.com/lithammer/shortuuid/v4"
	"gorm.io/gorm"
)

var (
	ErrInvalidBasePath    = errors.New("invalid base path")
	ErrNoSuchFile         = errors.New("no such file")
	ErrFileNotUploaded    = errors.New("file not uploaded")
	ErrInvalidMimeType    = errors.New("invalid mimetype")
	ErrUploadSizeExceeded = errors.New("upload size exceeded")
)

var AllowedMimeTypes = []string{
	"image/jpeg",
	"image/png",
	"image/gif",
	"image/webp",
}

const MimeHeaderSize_b = 512
const UploadSizeLimit_b = 1_000_000

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
	fi, err := os.Stat(basePath)
	if err != nil {
		return nil, err
	} else if !fi.IsDir() {
		return nil, ErrInvalidBasePath
	}
	// __jm__ make sure is writable?

	return &FileStorageServiceImpl{
		basePath:     basePath,
		baseServices: baseServices,
	}, nil
}

func (service *FileStorageServiceImpl) CreateStub(user *User) (*FileMetadata, error) {
	fileMetadata := &FileMetadata{
		PublicId: shortuuid.New(),
		OwnerId:  user.ID,
	}
	tx := service.baseServices.Database.Create(fileMetadata)
	if err := tx.GetError(); err != nil {
		return nil, err
	}
	return fileMetadata, nil
}

func (service *FileStorageServiceImpl) GetData(id string) (*os.File, error) {
	tx := service.baseServices.Database.Find(nil, FileMetadata{PublicId: id})
	if err := tx.GetError(); errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNoSuchFile // TODO: using errors.Join?
	} else if err != nil {
		return nil, err
	}

	file, err := os.Open(path.Join(service.basePath, id))
	if err != nil {
		return nil, ErrFileNotUploaded
	}

	return file, nil
}

// limit upload to ~1mb
func (service *FileStorageServiceImpl) Upload(fileMetadata FileMetadata, data io.Reader, mimetype string) (*FileMetadata, error) {
	// check mimetype is in allowedMimetypes slice
	dataBytes := []byte{}
	n, err := data.Read(dataBytes)
	if err != nil {
		return nil, err // __jm__ failed to read
	}
	if n > UploadSizeLimit_b {
		return nil, ErrUploadSizeExceeded
	}

	actualMimeType := http.DetectContentType(dataBytes)
	if actualMimeType != mimetype {
		return nil, ErrInvalidMimeType
	}

	file, err := os.Open(path.Join(service.basePath, fileMetadata.PublicId))
	if err != nil {
		return nil, err // __jm__ failed to open
	}

	if _, err = file.Write(dataBytes); err != nil {
		return nil, err // __jm__ failed to write
	}

	fileMetadata.ContentType = sql.NullString{String: mimetype, Valid: true}
	fileMetadata.Uploaded = sql.NullTime{Time: time.Now(), Valid: true}
	tx := service.baseServices.Database.Save(&fileMetadata)
	if err = tx.GetError(); err != nil {
		return nil, err // __jm__ failed to save
	}

	return &fileMetadata, nil
}

func (service *FileStorageServiceImpl) Remove(fileMetadata FileMetadata) error {
	err := os.Remove(path.Join(service.basePath, fileMetadata.PublicId))
	if err != nil {
		return ErrFileNotUploaded
	}

	return nil
}
