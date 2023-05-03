package middleware

import (
	"log"
	"net/http/httptest"
	"testing"

	api "github.com/StampWallet/backend/internal/api/models"
	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services/mocks"
	. "github.com/StampWallet/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func getAuthMiddleware(ctrl *gomock.Controller) *AuthMiddleware {
	return &AuthMiddleware{
		logger:       log.Default(),
		tokenService: NewMockTokenService(ctrl),
	}
}

// Q: Middleware should only look at Authorization header and assign user based on token value?
//    Should also look at user credentials, but nothing beyond that and future tests should reflect that

func TestHandleOk(t *testing.T) {
	// data prep
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	testToken := "ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK"
	testTokenId := "0123456789"
	testUser := GetDefaultUser()
	testTokenStruct := &Token{
		OwnerId:   testUser.ID,
		TokenId:   testTokenId,
		TokenHash: testToken,
		// Expires: _, TODO: needed? time now + 24h or custom rule
		TokenPurpose: TokenPurposeSession,
		Used:         false,
		Recalled:     false,
	}

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards").
		SetMethod("GET").
		// SetDefaultUser(). // Q: user should be part of request, but result of calling Handle should be a set "user" env var
		SetHeader("Authorization", "Bearer "+testToken).
		SetHeader("Accept", "application/json").
		Context
	// Q: token also in payload? doesnt make sense

	// test env prep
	ctrl := gomock.NewController(t)
	authMiddleware := getAuthMiddleware(ctrl)

	// TODO: mock accessor to token
	// accessor should expect token and return token id

	authMiddleware.tokenService.(*MockTokenService).
		EXPECT().
		Check(
			gomock.Eq(testTokenId),
			gomock.Eq(testToken),
		).
		Return(
			testUser,
			testTokenStruct,
			nil,
		)

	authMiddleware.Handle(context)

	user, ok := context.Get("user")
	if ok == false {
		t.Errorf("Context holds no user")
	}
	require.Truef(t, MatchEntities(testUser, user), "User entity provided by authenticator failed validation")
	// TODO: test MatchEntities usage
}

func TestHandleNok_UnknownToken(t *testing.T) {
	// data prep
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	testToken := "ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK"

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards").
		SetMethod("GET").
		SetHeader("Authorization", "Bearer "+testToken).
		SetHeader("Accept", "application/json").
		Context
	// Q: User creds w/ request
	// Q: token also in payload? doesnt make sense

	// test env prep
	ctrl := gomock.NewController(t)
	authMiddleware := getAuthMiddleware(ctrl)

	// TODO: mock accessor to token
	// accessor should expect token and return error saying token not found in db

	authMiddleware.Handle(context)

	respBody, respCode, err := ExtractResponse[api.DefaultResponse](t, w)

	require.Nilf(t, err, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(401), "Response returned unexpected status code")
	require.Equalf(t, respBody.Status, api.FORBIDDEN, "Status inside default response body does not match expected")

	userPtr, userExists := context.Get("user")
	require.Truef(t, userPtr == nil && userExists == false, "User field was overwritten despite no valid user existing")
}

func TestHandleNok_InvalidToken(t *testing.T) {
	// data prep
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	testToken := "ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK"
	testTokenId := "0123456789"

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/user/cards").
		SetMethod("GET").
		SetHeader("Authorization", "Bearer "+testToken).
		SetHeader("Accept", "application/json").
		Context
	// Q: token also in payload? doesnt make sense

	// test env prep
	ctrl := gomock.NewController(t)
	authMiddleware := getAuthMiddleware(ctrl)

	// TODO: mock accessor to token
	// accessor should expect token and return id

	authMiddleware.tokenService.(*MockTokenService).
		EXPECT().
		Check(
			gomock.Eq(testTokenId),
			gomock.Eq(testToken),
		).
		Return(
		// Q: returns what? token exists in db, but assigned to another user - return valid data or err out in tokenService?
		)

	authMiddleware.Handle(context)

	respBody, respCode, err := ExtractResponse[api.DefaultResponse](t, w)

	require.Nilf(t, err, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(401), "Response returned unexpected status code")
	require.Equalf(t, respBody.Status, api.FORBIDDEN, "Status inside default response body does not match expected")

	userPtr, userExists := context.Get("user")
	require.Truef(t, userPtr == nil && userExists == false, "User field was overwritten despite no valid user existing")
}
