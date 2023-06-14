package services

import (
	"database/sql"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"time"

	. "github.com/StampWallet/backend/internal/database"
	"github.com/StampWallet/backend/internal/utils"
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

// limit upload to ~1mb
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
	tx := service.baseServices.Database.Find(&FileMetadata{}, FileMetadata{PublicId: id})
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

func (service *FileStorageServiceImpl) Upload(fileMetadata FileMetadata, data io.Reader, mimetype string) (*FileMetadata, error) {
	dataBytes := []byte{}
	for {
		buf := make([]byte, 512)
		_, err := data.Read(buf)

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if len(dataBytes)+len(buf) > UploadSizeLimit_b {
			return nil, ErrUploadSizeExceeded
		}
		dataBytes = append(dataBytes, buf...)
	}

	actualMimeType := http.DetectContentType(dataBytes)
	if actualMimeType != mimetype || !utils.Contains(AllowedMimeTypes, actualMimeType) {
		return nil, ErrInvalidMimeType
	}

	err := os.WriteFile(
		path.Join(service.basePath, fileMetadata.PublicId),
		dataBytes,
		fs.FileMode(os.O_WRONLY),
	)
	if err != nil {
		return nil, err
	}

	fileMetadata.ContentType = sql.NullString{String: mimetype, Valid: true}
	fileMetadata.Uploaded = sql.NullTime{Time: time.Now().Round(time.Microsecond), Valid: true}
	tx := service.baseServices.Database.Save(&fileMetadata)
	if err = tx.GetError(); err != nil {
		return nil, err
	}

	return &fileMetadata, nil
}

func (service *FileStorageServiceImpl) Remove(fileMetadata FileMetadata) error {
	err := os.Remove(path.Join(service.basePath, fileMetadata.PublicId))
	if err != nil {
		return ErrFileNotUploaded
	}

	fileMetadata.ContentType = sql.NullString{}
	fileMetadata.Uploaded = sql.NullTime{}
	tx := service.baseServices.Database.Save(&fileMetadata)
	if err := tx.GetError(); err != nil {
		return err
	}

	return nil
}
