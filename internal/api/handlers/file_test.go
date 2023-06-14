package api

import (
	"bytes"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	api "github.com/StampWallet/backend/internal/api/models"
	"github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/database/accessors/mocks"
	. "github.com/StampWallet/backend/internal/services/mocks"
	. "github.com/StampWallet/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func getFileHandlers(ctrl *gomock.Controller) *FileHandlers {
	return &FileHandlers{
		fileStorageService:    NewMockFileStorageService(ctrl),
		logger:                log.Default(),
		userAuthorizedAcessor: NewMockUserAuthorizedAccessor(ctrl),
	}
}

func TestFileHandlerGetFileOk(t *testing.T) {
	testUser := GetDefaultUser()
	fileId := "abcdef123"

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/file/"+fileId).
		SetUser(testUser).
		SetMethod("GET").
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "multipart/form-data").
		SetParam("fileId", fileId).
		SetDefaultToken().
		Context

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getFileHandlers(ctrl)

	// setup mocks

	testFileHandle, err := os.Open("resources/test.png")
	require.NoError(t, err)

	handler.fileStorageService.(*MockFileStorageService).
		EXPECT().
		GetData(gomock.Eq(fileId)).
		Return(
			testFileHandle, // Q: Should this change upon upload?
			"image/png",
			nil,
		)

	handler.getFile(context)

	expectedContents, _ := io.ReadAll(TestFileReader("resources/test.png"))

	require.Equalf(t, w.Result().StatusCode, int(200), "Response returned unexpected status code")
	require.Truef(t, bytes.Compare(w.Body.Bytes(), expectedContents) == 0, "Response returned unexpected file data")
}

func TestFileHandlersPostFileOk(t *testing.T) {
	testUser := GetDefaultUser()
	fileId := "abcdef123"
	testFileMetadata := GetTestFileMetadata(nil, testUser)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/file/"+fileId).
		SetUser(testUser).
		SetMethod("POST").
		AttachTestPng().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "image/png").
		SetParam("fileId", fileId).
		SetDefaultToken().
		Context

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getFileHandlers(ctrl)

	// setup mocks
	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		Get(testUser, &database.FileMetadata{PublicId: fileId}).
		Return(testFileMetadata, nil)

	handler.fileStorageService.(*MockFileStorageService).
		EXPECT().
		Upload(
			gomock.Eq(*testFileMetadata),
			gomock.Any(), // TODO: Matcher for test png?
			gomock.Eq("image/png"),
		).
		Return(
			testFileMetadata, // Q: Should this change upon upload?
			nil,
		)

	handler.postFile(context)

	respBodyExpected := api.DefaultResponse{Status: api.OK}
	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, reflect.DeepEqual(respBodyExpected, *respBody), "Response returned unexpected body contents")
	// TODO: test MatchEntities and gomock.Eq
}

func TestFileHandlersDeleteFileOk(t *testing.T) {
	testUser := GetDefaultUser()
	fileId := "abcdef123"
	testFileMetadata := GetTestFileMetadata(nil, testUser)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/file/"+fileId).
		SetUser(testUser).
		SetMethod("DELETE").
		SetHeader("Accept", "application/json").
		SetParam("fileId", fileId).
		SetDefaultToken().
		Context

	// test env prep
	ctrl := gomock.NewController(t)
	handler := getFileHandlers(ctrl)

	handler.userAuthorizedAcessor.(*MockUserAuthorizedAccessor).
		EXPECT().
		Get(testUser, &database.FileMetadata{PublicId: fileId}).
		Return(testFileMetadata, nil)

	handler.fileStorageService.(*MockFileStorageService).
		EXPECT().
		RemoveFile(gomock.Eq(*testFileMetadata)).
		Return(nil)

	handler.deleteFile(context)

	respBodyExpected := api.DefaultResponse{Status: api.OK}
	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, reflect.DeepEqual(respBodyExpected, *respBody), "Response returned unexpected body contents")
}
