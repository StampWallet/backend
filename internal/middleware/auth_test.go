package middleware

import (
	"log"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	api "github.com/StampWallet/backend/internal/api/models"
	. "github.com/StampWallet/backend/internal/database"
	"github.com/StampWallet/backend/internal/services"
	. "github.com/StampWallet/backend/internal/services/mocks"
	. "github.com/StampWallet/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// Create AuthMiddleware
func getAuthMiddleware(ctrl *gomock.Controller) *AuthMiddleware {
	return &AuthMiddleware{
		logger:       log.Default(),
		tokenService: NewMockTokenService(ctrl),
	}
}

// Test AuthMiddleware on the happy path
func TestHandleOk(t *testing.T) {
	// data prep
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	testTokenSecret := "ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK"
	testTokenId := "0123456789"
	testToken := testTokenId + ":" + testTokenSecret
	testUser := GetDefaultUser()
	testTokenStruct := &Token{
		OwnerId:      testUser.ID,
		TokenId:      testTokenId,
		TokenHash:    testToken,
		Expires:      time.Now().Add(time.Hour * 24),
		TokenPurpose: TokenPurposeSession,
		Used:         false,
		Recalled:     false,
		User:         testUser,
	}

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards").
		SetMethod("GET").
		SetHeader("Authorization", "Bearer "+testToken).
		SetHeader("Accept", "application/json").
		Context

	// test env prep
	ctrl := gomock.NewController(t)
	authMiddleware := getAuthMiddleware(ctrl)

	authMiddleware.tokenService.(*MockTokenService).
		EXPECT().
		Check(
			gomock.Eq(testTokenId),
			gomock.Eq(testTokenSecret),
		).
		Return(
			testTokenStruct,
			nil,
		)

	authMiddleware.Handle(context)

	user, ok := context.Get("user")
	if ok == false {
		t.Errorf("Context holds no user")
	}
	require.Truef(t, reflect.DeepEqual(testUser, user), "User entity provided by authenticator failed validation")
	// TODO: test MatchEntities usage
}

func TestHandleNok_EmailToken(t *testing.T) {
	// data prep
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	testTokenSecret := "ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK"
	testTokenId := "0123456789"
	testToken := testTokenId + ":" + testTokenSecret
	testUser := GetDefaultUser()
	testTokenStruct := &Token{
		OwnerId:      testUser.ID,
		TokenId:      testTokenId,
		TokenHash:    testToken,
		Expires:      time.Now().Add(time.Hour * 24),
		TokenPurpose: TokenPurposeEmail,
		Used:         false,
		Recalled:     false,
		User:         testUser,
	}

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards").
		SetMethod("GET").
		SetHeader("Authorization", "Bearer "+testToken).
		SetHeader("Accept", "application/json").
		Context

	// test env prep
	ctrl := gomock.NewController(t)
	authMiddleware := getAuthMiddleware(ctrl)

	authMiddleware.tokenService.(*MockTokenService).
		EXPECT().
		Check(
			gomock.Eq(testTokenId),
			gomock.Eq(testTokenSecret),
		).
		Return(
			testTokenStruct,
			nil,
		)

	authMiddleware.Handle(context)

	respBodyExpected := api.DefaultResponse{Status: api.UNAUTHORIZED}
	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(401), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Status inside default response body does not match expected")
	// TODO: MatchEntities

	userPtr, userExists := context.Get("user")
	require.Truef(t, userPtr == nil && userExists == false, "User field was overwritten despite no valid user existing")
}

// Test AuthMiddleware when the token is invalid
func TestHandleNok_UnknownToken(t *testing.T) {
	// data prep
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	testTokenSecret := "ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK"
	testTokenId := "0123456789"
	testToken := testTokenId + ":" + testTokenSecret

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards").
		SetMethod("GET").
		SetHeader("Authorization", "Bearer "+testToken).
		SetHeader("Accept", "application/json").
		Context

	// test env prep
	ctrl := gomock.NewController(t)
	authMiddleware := getAuthMiddleware(ctrl)

	authMiddleware.tokenService.(*MockTokenService).
		EXPECT().
		Check(
			gomock.Eq(testTokenId),
			gomock.Eq(testTokenSecret),
		).
		Return(
			nil,
			err,
		)

	authMiddleware.Handle(context)

	respBodyExpected := api.DefaultResponse{Status: api.UNAUTHORIZED}
	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(401), "Response returned unexpected status code")
	require.Truef(t, reflect.DeepEqual(respBodyExpected, *respBody), "Status inside default response body does not match expected")
	// TODO: MatchEntities

	userPtr, userExists := context.Get("user")
	require.Truef(t, userPtr == nil && userExists == false, "User field was overwritten despite no valid user existing")
}

func TestHandleNok_TokenErrors(t *testing.T) {
	for _, err := range []error{services.ErrUnknownToken, services.ErrTokenExpired, services.ErrTokenUsed} {
		testHandleTokenError(t, err)
	}
}
