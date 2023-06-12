package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lithammer/shortuuid/v4"
	"github.com/stretchr/testify/require"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/testutils"
)

func GetFileStorageService(ctrl *gomock.Controller) *FileStorageServiceImpl {
	temp, err := os.MkdirTemp("", "stampwallet_test")
	if err != nil {
		panic(fmt.Sprintf("failed to create temp dir %s", err))
	}
	return &FileStorageServiceImpl{
		basePath: temp,
		baseServices: BaseServices{
			Logger:   log.Default(),
			Database: GetTestDatabase(),
		},
	}
}

func createFile(t *testing.T, service *FileStorageServiceImpl, publicId string) (*os.File, string) {
	file, err := os.Create(path.Join(service.basePath, publicId))
	require.Nilf(t, err, "os.Create returned an error")

	toWrite := shortuuid.New()
	n, err := file.Write([]byte(toWrite))
	require.Nilf(t, err, "file.Write returned an error")
	require.Equalf(t, len(toWrite), n, "file.Write returned an error")

	return file, toWrite
}

func readAndCompare(t *testing.T, org string, service *FileStorageServiceImpl, publicId string) {
	newFile, err := service.GetData(publicId)
	require.Nilf(t, err, "service.GetData returned an error")

	toRead := make([]byte, len(org))
	newFile.Read(toRead)
	require.Equalf(t, org, string(toRead), "file.Read returned a different file")
}

// Tests

func TestFileStorageServiceCreateStub(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := GetFileStorageService(ctrl)
	defer os.RemoveAll(service.basePath)
	user := GetTestUser(service.baseServices.Database)

	metadata, err := service.CreateStub(user)
	require.Nilf(t, err, "Error should be nil")
	require.NotNilf(t, metadata, "FileMetadata should not be nil")
	require.Equalf(t, user.ID, metadata.OwnerId, "Metadata has invalid owner")
	require.Falsef(t, metadata.Uploaded.Valid, "Metadata has upload date")
}

func TestFileStorageServiceGetData(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := GetFileStorageService(ctrl)
	defer os.RemoveAll(service.basePath)
	user := GetTestUser(service.baseServices.Database)

	metadata := FileMetadata{
		PublicId:    shortuuid.New(),
		OwnerId:     user.ID,
		Uploaded:    sql.NullTime{Time: time.Now(), Valid: true},
		ContentType: sql.NullString{String: AllowedMimeTypes[0], Valid: true},
	}
	tx := service.baseServices.Database.Create(&metadata)
	require.Nilf(t, tx.GetError(), "Database.Create returned an error")

	_, toWrite := createFile(t, service, metadata.PublicId)
	readAndCompare(t, toWrite, service, metadata.PublicId)
}

func TestFileStorageServiceGetDataInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := GetFileStorageService(ctrl)
	defer os.RemoveAll(service.basePath)

	file, err := service.GetData("invalid uuid lol")
	require.Nilf(t, file, "service.GetData returned a file")
	require.ErrorAsf(t, err, &ErrNoSuchFile, "service.GetData returned a file")
}

func TestFileStorageServiceGetDataNotUploaded(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := GetFileStorageService(ctrl)
	defer os.RemoveAll(service.basePath)
	user := GetTestUser(service.baseServices.Database)

	metadata := FileMetadata{
		PublicId: shortuuid.New(),
		OwnerId:  user.ID,
		Uploaded: sql.NullTime{Valid: false},
	}
	tx := service.baseServices.Database.Create(&metadata)
	require.Nilf(t, tx.GetError(), "Database.Create returned an error")

	newFile, err := service.GetData(metadata.PublicId)
	require.Nilf(t, newFile, "service.GetData returned a file")
	require.ErrorAsf(t, err, &ErrFileNotUploaded, "service.GetData returned a file")
}

func TestFileStorageServiceUpload(t *testing.T) {
	// __jm__ refactor tests so data includes mimetypes

	ctrl := gomock.NewController(t)
	service := GetFileStorageService(ctrl)
	defer os.RemoveAll(service.basePath)
	user := GetTestUser(service.baseServices.Database)

	metadata := FileMetadata{
		PublicId: shortuuid.New(),
		OwnerId:  user.ID,
	}
	tx := service.baseServices.Database.Create(&metadata)
	require.Nilf(t, tx.GetError(), "Database.Create returned an error")

	_, err := os.Create(path.Join(service.basePath, metadata.PublicId))
	require.Nilf(t, err, "os.Create returned an error")

	file, toWrite := createFile(t, service, shortuuid.New())
	_, err = file.Seek(0, 0)
	require.Nilf(t, err, "file.Seek returned an error")

	newFileMetadata, err := service.Upload(metadata, file, AllowedMimeTypes[0])
	require.Nilf(t, err, "service.GetData returned an error")
	require.NotNilf(t, newFileMetadata, "service.GetData returned a nil FileMetadata")
	require.Truef(t, newFileMetadata.Uploaded.Valid, "newFileMetadata has invalid upload time")
	require.Truef(t, TimeJustAroundNow(newFileMetadata.Uploaded.Time), "newFileMetadata has invalid upload time")
	require.Truef(t, newFileMetadata.ContentType.Valid, "newFileMetadata has invalid content type")
	require.Equalf(t, AllowedMimeTypes[0], newFileMetadata.ContentType.String,
		"newFileMetadata has invalid content type")

	readAndCompare(t, toWrite, service, metadata.PublicId)

	var fileMetadataDb FileMetadata
	tx = service.baseServices.Database.First(&fileMetadataDb,
		FileMetadata{PublicId: newFileMetadata.PublicId})
	require.Nilf(t, tx.GetError(), "Database.First returned an error")
	require.Equal(t, newFileMetadata.Uploaded.Time, fileMetadataDb.Uploaded.Time)
	require.Equal(t, newFileMetadata.ContentType.String, fileMetadataDb.ContentType.String)
	require.Equal(t, newFileMetadata.OwnerId, fileMetadataDb.OwnerId)
}

func TestFileStorageServiceUploadInvalidMimeType(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := GetFileStorageService(ctrl)
	defer os.RemoveAll(service.basePath)
	user := GetTestUser(service.baseServices.Database)

	metadata := FileMetadata{
		PublicId: shortuuid.New(),
		OwnerId:  user.ID,
	}
	tx := service.baseServices.Database.Create(&metadata)
	require.Nilf(t, tx.GetError(), "Database.Create returned an error")

	_, err := os.Create(path.Join(service.basePath, metadata.PublicId))
	require.Nilf(t, err, "os.Create returned an error")

	file, _ := createFile(t, service, shortuuid.New())

	newFileMetadata, err := service.Upload(metadata, file, "test/test")
	require.ErrorAsf(t, err, ErrInvalidMimeType, "FileStorageService.Upload did not return InvalidMimeType")
	require.Nilf(t, newFileMetadata, "FileStorageService.Upload did not return nil FileMetadata")
}

func TestFileStorageServiceRemove(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := GetFileStorageService(ctrl)
	defer os.RemoveAll(service.basePath)
	user := GetTestUser(service.baseServices.Database)

	metadata := FileMetadata{
		PublicId:    shortuuid.New(),
		OwnerId:     user.ID,
		Uploaded:    sql.NullTime{Time: time.Now(), Valid: true},
		ContentType: sql.NullString{String: AllowedMimeTypes[0], Valid: true},
	}
	tx := service.baseServices.Database.Create(&metadata)
	require.Nilf(t, tx.GetError(), "Database.Create returned an error")

	file, _ := createFile(t, service, metadata.PublicId)
	file.Close()

	err := service.Remove(metadata)
	require.Nilf(t, err, "FileStorageService.Remove returned an error")

	openedFile, err := os.Open(path.Join(service.basePath, metadata.PublicId))
	require.Nilf(t, openedFile, "os.Open returned a file - file exists, but should have been removed")
	require.Errorf(t, err, "os.Open did not return an error")

	serviceFile, err := service.GetData(metadata.PublicId)
	require.Nilf(t, serviceFile, "service.GetData returned a file")
	require.ErrorAs(t, err, ErrFileNotUploaded, "service.GetData did not return a FileNotUploaded error")

	var fileMetadataDb FileMetadata
	tx = service.baseServices.Database.First(&fileMetadataDb,
		FileMetadata{PublicId: metadata.PublicId})
	require.Nilf(t, tx.GetError(), "Database.First returned an error")
	require.Falsef(t, fileMetadataDb.Uploaded.Valid, "fileMetadataDb.Uploaded.Valid is not false")
	require.Falsef(t, fileMetadataDb.ContentType.Valid, "fileMetadataDb.ContentType.Valid is not false")
}

func TestFileStorageServiceRemoveNotUploaded(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := GetFileStorageService(ctrl)
	defer os.RemoveAll(service.basePath)
	user := GetTestUser(service.baseServices.Database)

	metadata := FileMetadata{
		PublicId:    shortuuid.New(),
		OwnerId:     user.ID,
		Uploaded:    sql.NullTime{Valid: false},
		ContentType: sql.NullString{Valid: false},
	}
	tx := service.baseServices.Database.Create(&metadata)
	require.Nilf(t, tx.GetError(), "Database.Create returned an error")

	err := service.Remove(metadata)
	require.ErrorAsf(t, err, &ErrFileNotUploaded, "FileStorageService.Remove did not return FileNotUploaded")

	var fileMetadataDb FileMetadata
	tx = service.baseServices.Database.First(&fileMetadataDb,
		FileMetadata{PublicId: metadata.PublicId})
	require.Nilf(t, tx.GetError(), "Database.First returned an error")
	require.Falsef(t, fileMetadataDb.Uploaded.Valid, "fileMetadataDb.Uploaded.Valid is not false")
	require.Falsef(t, fileMetadataDb.ContentType.Valid, "fileMetadataDb.ContentType.Valid is not false")
}
