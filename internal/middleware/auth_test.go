package middleware

import (
	"log"
	"net/http/httptest"
	"testing"
	"time"

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

	testToken := "0123456789:ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK"

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

	authMiddleware.Handle(context)

	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](t, w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(401), "Response returned unexpected status code")
	require.Equalf(t, respBody.Status, api.FORBIDDEN, "Status inside default response body does not match expected")

	userPtr, userExists := context.Get("user")
	require.Truef(t, userPtr == nil && userExists == false, "User field was overwritten despite no valid user existing")
}
