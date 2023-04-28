package api

import (
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	api "github.com/StampWallet/backend/internal/api/models"
	"github.com/StampWallet/backend/internal/database"
	"github.com/StampWallet/backend/internal/managers"
	. "github.com/StampWallet/backend/internal/managers/mocks"
	. "github.com/StampWallet/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func GetAuthHandlers(ctrl *gomock.Controller) *AuthHandlers {
	return &AuthHandlers{
		authManager: NewMockAuthManager(ctrl),
		logger:      log.Default(),
	}
}

func TestAuthHandlersPostAccount(t *testing.T) {
	// data prep
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	payload := api.PostAccountRequest{
		Email:    "test@example.com",
		Password: "zaq1@WSX",
	}
	payloadJson, _ := json.Marshal(payload)

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/auth/account").
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(payloadJson).
		Context

	userDetailsExpected := managers.UserDetails{
		Email:    "test@example.com",
		Password: "zaq1@WSX",
	}

	respBodyExpected := api.DefaultResponse{Status: api.OK}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		Create(gomock.Eq(userDetailsExpected))

	handler.postAccount(context)

	respBody, respCode, err := ExtractResponse[api.DefaultResponse](t, w)

	require.Nilf(t, err, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(201), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBody, respBodyExpected), "Response returned unexpected body contents")
}

func TestAuthHandlersPostAccountEmail(t *testing.T) {
	// data prep
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	payload := api.PostAccountEmailRequest{
		Email: "test@example.com",
	}
	payloadJson, _ := json.Marshal(payload)

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/auth/account/email").
		SetDefaultUser().
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(payloadJson).
		Context

	userExpected := GetDefaultUser()
	newEmailExpected := "test@example.com"
	respBodyExpected := api.DefaultResponse{Status: api.OK}

	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		ChangeEmail(
			gomock.Eq(userExpected),
			gomock.Eq(newEmailExpected),
		).
		Return(
			userExpected,
			nil,
		)

	handler.postAccountEmail(context)

	respBody, respCode, err := ExtractResponse[api.DefaultResponse](t, w)

	require.Nilf(t, err, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBody, respBodyExpected), "Response returned unexpected body contents")
}

func TestAuthHandlersPostAccountPasswordOk(t *testing.T) {
	// data prep
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	payload := api.PostAccountPasswordRequest{
		OldPassword: "zaq1@WSX",
		Password:    "XSW@1qaz",
	}
	payloadJson, _ := json.Marshal(payload)

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/auth/account/password").
		SetDefaultUser().
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(payloadJson).
		Context

	userExpected := GetDefaultUser()
	oldPasswordExpected := "zaq1@WSX"
	newPasswordExpected := "XSW@1qaz"
	respBodyExpected := api.DefaultResponse{Status: api.OK}

	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		ChangePassword(
			gomock.Eq(userExpected),
			gomock.Eq(oldPasswordExpected),
			gomock.Eq(newPasswordExpected),
		).
		Return(
			userExpected,
			nil,
		)

	handler.postAccountPassword(context)

	respBody, respCode, err := ExtractResponse[api.DefaultResponse](t, w)

	require.Nilf(t, err, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBody, respBodyExpected), "Response returned unexpected body contents")
}

func TestAuthHandlersPostAccountPasswordNok_old_pass(t *testing.T) {
	// data prep
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	payload := api.PostAccountPasswordRequest{
		OldPassword: "zaq1@WSX",
		Password:    "zaq1@WSX",
	}
	payloadJson, _ := json.Marshal(payload)

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/auth/account/password").
		SetDefaultUser().
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(payloadJson).
		Context

	respBodyExpected := api.DefaultResponse{Status: api.INVALID_REQUEST}

	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.postAccountPassword(context)

	respBody, respCode, err := ExtractResponse[api.DefaultResponse](t, w)

	require.Nilf(t, err, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(400), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBody, respBodyExpected), "Response returned unexpected body contents")
}

// TODO
func TestAuthHandlersPostAccountEmailConfirmation(t *testing.T) {
	// data prep
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	exampleToken := "ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK"
	payload := api.PostAccountEmailConfirmationRequest{
		Token: exampleToken,
	}
	payloadJson, _ := json.Marshal(payload)

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/auth/account/emailConfirmation").
		SetDefaultUser(). // Authorization already performed by middleware
		SetMethod("POST").
		SetHeader("Authorization", "Bearer "+exampleToken).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(payloadJson).
		Context

	userExpected := GetDefaultUser()
	respBodyExpected := api.DefaultResponse{Status: api.OK}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		ConfirmEmail(
			gomock.Eq(userExpected), //TODO Correct arguments and generate new mocks
			gomock.Eq(exampleToken),
		).Return(nil)

	handler.postAccountEmailConfirmation(context)

	respBody, respCode, err := ExtractResponse[api.DefaultResponse](t, w)

	require.Nilf(t, err, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBody, respBodyExpected), "Response returned unexpected body contents")
}

func TestAuthHandlersPostSession(t *testing.T) {
	// data prep
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	payload := api.PostAccountSessionRequest{
		Email:    "test@example.com",
		Password: "zaq1@WSX",
	}
	payloadJson, _ := json.Marshal(payload)

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/auth/account/sessions").
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(payloadJson).
		Context

	respBodyExpected := api.PostAccountSessionResponse{
		Token: "ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK",
	}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		Login(
			gomock.Eq("test@example.com"),
			gomock.Eq("zaq1@WSX"),
		).
		Return(
			&database.User{
				PublicId:      "testUserId",
				Email:         "test@example.com",
				PasswordHash:  "hash of zaq1@WSX",
				EmailVerified: true,
			},
			&database.Token{
				OwnerId:   0,
				TokenId:   "testTokenId",
				TokenHash: "ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK",
				// Expires: _, TODO - see how to mock this
				TokenPurpose: database.TokenPurposeSession,
				Used:         false,
				Recalled:     false,
			},
			nil,
		)

	handler.postSession(context)

	respBody, respCode, err := ExtractResponse[api.PostAccountSessionResponse](t, w)

	require.Nilf(t, err, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(201), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBody, respBodyExpected), "Response returned unexpected body contents")
}

func TestAuthHandlersDeleteSession(t *testing.T) {
	// data prep
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/auth/sessions").
		SetDefaultUser().
		SetMethod("DELETE").
		SetHeader("Authorization", "Bearer ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK").
		SetHeader("Accept", "application/json").
		Context

	respBodyExpected := api.DefaultResponse{Status: api.OK}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		Logout(
			gomock.Eq("test@example.com"),
			gomock.Eq("zaq1@WSX"),
		).
		Return(
			&database.User{
				PublicId:      "testUserId",
				Email:         "test@example.com",
				PasswordHash:  "hash of zaq1@WSX",
				EmailVerified: true,
			},
			&database.Token{
				OwnerId:   0,
				TokenId:   "testTokenId",
				TokenHash: "ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK",
				// Expires: _, TODO - see how F mocks this
				TokenPurpose: database.TokenPurposeSession,
				Used:         false,
				Recalled:     false,
			},
			nil,
		)

	handler.postSession(context)

	respBody, respCode, err := ExtractResponse[api.PostAccountSessionResponse](t, w)

	require.Nilf(t, err, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBody, respBodyExpected), "Response returned unexpected body contents")
}
