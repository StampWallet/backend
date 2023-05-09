package api

import (
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	api "github.com/StampWallet/backend/internal/api/models"
	"github.com/StampWallet/backend/internal/database"
	"github.com/StampWallet/backend/internal/managers"
	. "github.com/StampWallet/backend/internal/managers/mocks"
	. "github.com/StampWallet/backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func GetAuthHandlers(ctrl *gomock.Controller) *AuthHandlers {
	return &AuthHandlers{
		authManager: NewMockAuthManager(ctrl),
		logger:      log.Default(),
	}
}

func SetupAuthHandlersPostAccount() (
	w *httptest.ResponseRecorder,
	context *gin.Context,
	userDetailsExpectedPtr *managers.UserDetails,
	testUser *database.User,
	testToken *database.Token,
) {
	testEmail := "test@example.com"
	testPassword := "zaq1@WSX"

	testUser = GetDefaultUser()
	testUser.Email = testEmail
	hash, _ := bcrypt.GenerateFromPassword([]byte(testPassword), 10)
	testUser.PasswordHash = string(hash)
	testUser.EmailVerified = false

	testToken = &database.Token{
		OwnerId:      testUser.ID,
		TokenId:      "01234556789",
		TokenHash:    "u8m932r98u3", // TODO: more fitting test value
		Expires:      time.Now().Add(time.Hour * 24),
		TokenPurpose: database.TokenPurposeEmail,
		Used:         false,
		Recalled:     false,
	}

	payload := api.PostAccountRequest{
		Email:    testEmail,
		Password: testPassword,
	}
	payloadJson, _ := json.Marshal(payload)

	gin.SetMode(gin.TestMode)
	w = httptest.NewRecorder()

	context = NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/auth/account").
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(payloadJson).
		Context

	userDetailsExpected := managers.UserDetails{
		Email:    testEmail,
		Password: testPassword,
	}

	return w, context, &userDetailsExpected, testUser, testToken
}

func TestAuthHandlersPostAccountOk(t *testing.T) {
	w, context, userDetailsExpectedPtr, testUser, testToken := SetupAuthHandlersPostAccount()
	respBodyExpected := api.DefaultResponse{Status: api.OK}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		Create(gomock.Eq(*userDetailsExpectedPtr)).
		Return(
			testUser,
			testToken,
			nil,
		)

	handler.postAccount(context)

	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(201), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
	// TODO: test MatchEntities and gomock.Eq
}

func TestAuthHandlersPostAccountNok_DupMail(t *testing.T) {
	w, context, userDetailsExpectedPtr, _, _ := SetupAuthHandlersPostAccount()
	respBodyExpected := api.DefaultResponse{Status: api.CONFLICT}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		Create(gomock.Eq(*userDetailsExpectedPtr)).
		Return(
			nil,
			nil,
			managers.EmailExists,
		)

	handler.postAccount(context)

	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(409), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
	// TODO: test MatchEntities and gomock.Eq
}

func SetupAuthHandlersPostAccountEmail() (
	w *httptest.ResponseRecorder,
	testUser *database.User,
	testNewEmailPtr *string,
	context *gin.Context,
	testUserChangedEmail *database.User,
) {
	// data prep
	gin.SetMode(gin.TestMode)
	w = httptest.NewRecorder()

	testUser = GetDefaultUser()
	testNewEmail := "new_email_test@example.com"
	payload := api.PostAccountEmailRequest{
		Email: testNewEmail,
	}
	payloadJson, _ := json.Marshal(payload)

	context = NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/auth/account/email").
		SetUser(testUser).
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		SetBody(payloadJson).
		Context

	testUserChangedEmail = testUser
	testUserChangedEmail.Email = testNewEmail
	testUserChangedEmail.EmailVerified = false

	return w, testUser, &testNewEmail, context, testUserChangedEmail
}

func TestAuthHandlersPostAccountEmailOk(t *testing.T) {
	w, testUser, testNewEmailPtr, context, testUserChangedEmail := SetupAuthHandlersPostAccountEmail()
	respBodyExpected := api.DefaultResponse{Status: api.OK}

	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		ChangeEmail(
			gomock.Eq(testUser),
			gomock.Eq(&testNewEmailPtr),
		).
		Return(
			testUserChangedEmail,
			nil,
		)

	handler.postAccountEmail(context)

	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
	// TODO: MatchEntities and gomock.Eq
}

func TestAuthHandlersPostAccountEmailNok_DupMail(t *testing.T) {
	w, testUser, testNewEmailPtr, context, _ := SetupAuthHandlersPostAccountEmail()
	respBodyExpected := api.DefaultResponse{Status: api.CONFLICT}

	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		ChangeEmail(
			gomock.Eq(testUser),
			gomock.Eq(&testNewEmailPtr),
		).
		Return(
			nil,
			managers.EmailExists,
		)

	handler.postAccountEmail(context)

	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(409), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
	// TODO: MatchEntities and gomock.Eq
}

func TestAuthHandlersPostAccountPasswordOk(t *testing.T) {
	// data prep
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	oldPassword := "zaq1@WSX"
	newPassword := "XSW@1qaz"

	userOldPass := GetDefaultUser()
	hash, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), 10)
	userOldPass.PasswordHash = string(hash)

	userNewPass := userOldPass
	hash, _ = bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	userNewPass.PasswordHash = string(hash)

	payload := api.PostAccountPasswordRequest{
		OldPassword: oldPassword,
		Password:    newPassword,
	}
	payloadJson, _ := json.Marshal(payload)

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/auth/account/password").
		SetUser(userOldPass).
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		SetBody(payloadJson).
		Context

	respBodyExpected := api.DefaultResponse{Status: api.OK}

	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		ChangePassword(
			gomock.Eq(userOldPass),
			gomock.Eq(oldPassword),
			gomock.Eq(newPassword),
		).
		Return(
			userNewPass,
			nil,
		)

	handler.postAccountPassword(context)

	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
	// TODO: MatchEntities and gomock.Eq
}

func TestAuthHandlersPostAccountPasswordNok_OldPass(t *testing.T) {
	// data prep
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	testPassword := "zaq1@WSX"

	payload := api.PostAccountPasswordRequest{
		OldPassword: testPassword,
		Password:    testPassword,
	}
	payloadJson, _ := json.Marshal(payload)

	context := NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/auth/account/password").
		SetDefaultUser().
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetDefaultToken().
		SetBody(payloadJson).
		Context

	respBodyExpected := api.DefaultResponse{Status: api.INVALID_REQUEST}

	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.postAccountPassword(context)

	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(400), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBody, respBodyExpected), "Response returned unexpected body contents")
	// TODO: MatchEntities
}

func SetupAuthHandlersPostAccountEmailConfirmation() (
	w *httptest.ResponseRecorder,
	context *gin.Context,
	testUser *database.User,
	tokenId string,
	tokenSecret string,
) {
	// data prep
	gin.SetMode(gin.TestMode)
	w = httptest.NewRecorder()

	tokenId = "0123456789"
	tokenSecret = "ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK"
	exampleToken := tokenId + ":" + tokenSecret

	testUser = GetDefaultUser()

	payload := api.PostAccountEmailConfirmationRequest{
		Token: exampleToken,
	}
	payloadJson, _ := json.Marshal(payload)

	context = NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/auth/account/emailConfirmation").
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(payloadJson).
		Context

	return w, context, testUser, tokenId, tokenSecret
}

func TestAuthHandlersPostAccountEmailConfirmationOk(t *testing.T) {
	w, context, testUser, tokenId, tokenSecret := SetupAuthHandlersPostAccountEmailConfirmation()

	respBodyExpected := api.DefaultResponse{Status: api.OK}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		ConfirmEmail(
			gomock.Eq(tokenId),
			gomock.Eq(tokenSecret),
		).
		Return(
			testUser,
			nil,
		)

	handler.postAccountEmailConfirmation(context)

	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBody, respBodyExpected), "Response returned unexpected body contents")
	// TODO: MatchEntities and gomock.Eq
}

func TestAuthHandlersPostAccountEmailConfirmationNok_InvTok(t *testing.T) {
	w, context, _, tokenId, tokenSecret := SetupAuthHandlersPostAccountEmailConfirmation()

	respBodyExpected := api.DefaultResponse{Status: api.NOT_FOUND}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		ConfirmEmail(
			gomock.Eq(tokenId),
			gomock.Eq(tokenSecret),
		).
		Return(
			nil,
			managers.InvalidToken,
		)

	handler.postAccountEmailConfirmation(context)

	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(404), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBody, respBodyExpected), "Response returned unexpected body contents")
	// TODO: MatchEntities and gomock.Eq
}

func SetupAuthHandlersPostSession() (
	w *httptest.ResponseRecorder,
	context *gin.Context,
	testUser *database.User,
	testPassword string,
	testToken *database.Token,
	testTokenSecret string,
) {
	// data prep
	gin.SetMode(gin.TestMode)
	w = httptest.NewRecorder()

	testPassword = "zaq1@WSX"
	testUser = GetDefaultUser()
	hash, _ := bcrypt.GenerateFromPassword([]byte(testPassword), 10)
	testUser.PasswordHash = string(hash)

	payload := api.PostAccountSessionRequest{
		Email:    testUser.Email,
		Password: testPassword,
	}
	payloadJson, _ := json.Marshal(payload)

	context = NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/auth/account/sessions").
		SetMethod("POST").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(payloadJson).
		Context

	testTokenSecret = "ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK"
	hash, _ = bcrypt.GenerateFromPassword([]byte(testTokenSecret), 10)
	testToken = &database.Token{
		OwnerId:      testUser.ID,
		TokenId:      "testTokenId",
		TokenHash:    string(hash),
		Expires:      time.Now().Add(time.Hour * 24),
		TokenPurpose: database.TokenPurposeSession,
		Used:         false,
		Recalled:     false,
	}
	testUser.Tokens = append(testUser.Tokens, *testToken)

	return w, context, testUser, testPassword, testToken, testTokenSecret
}

func TestAuthHandlersPostSessionOk(t *testing.T) {
	w, context, testUser, testPassword, testToken, testTokenSecret := SetupAuthHandlersPostSession()

	// test env prep
	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		Login(
			gomock.Eq(testUser.Email),
			gomock.Eq(testPassword),
		).
		Return(
			testUser,
			testToken,
			nil,
		)

	handler.postSession(context)

	respBodyExpected := api.PostAccountSessionResponse{
		Token: testTokenSecret,
	}

	respBody, respCode, respParseErr := ExtractResponse[api.PostAccountSessionResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(201), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBody, respBodyExpected), "Response returned unexpected body contents")
	// TODO: MatchEntities and gomock.Eq
}

func TestAuthHandlersPostSessionNok_InvEmail(t *testing.T) {
	w, context, testUser, testPassword, _, _ := SetupAuthHandlersPostSession()

	// test env prep
	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		Login(
			gomock.Eq(testUser.Email),
			gomock.Eq(testPassword),
		).
		Return(
			nil,
			nil,
			managers.InvalidEmail,
		)

	handler.postSession(context)

	respBodyExpected := api.DefaultResponse{Status: api.UNAUTHORIZED}

	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(401), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBody, respBodyExpected), "Response returned unexpected body contents")
	// TODO: MatchEntities and gomock.Eq
}

func TestAuthHandlersPostSessionNok_InvLogin(t *testing.T) {
	w, context, testUser, testPassword, _, _ := SetupAuthHandlersPostSession()

	// test env prep
	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		Login(
			gomock.Eq(testUser.Email),
			gomock.Eq(testPassword),
		).
		Return(
			nil,
			nil,
			managers.InvalidLogin,
		)

	handler.postSession(context)

	respBodyExpected := api.DefaultResponse{Status: api.UNAUTHORIZED}

	respBody, respCode, respParseErr := ExtractResponse[api.DefaultResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(401), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBody, respBodyExpected), "Response returned unexpected body contents")
	// TODO: MatchEntities and gomock.Eq
}

func SetupAuthHandlersDeleteSession() (
	w *httptest.ResponseRecorder,
	context *gin.Context,
	testUser *database.User,
	tokenId string,
	tokenSecret string,
	testTokenStruct *database.Token,
) {
	// data prep
	gin.SetMode(gin.TestMode)
	w = httptest.NewRecorder()

	tokenId = "012345789"
	tokenSecret = "ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK"
	testToken := tokenId + ":" + tokenSecret

	testUser = GetDefaultUser()
	tokenHash, _ := bcrypt.GenerateFromPassword([]byte(tokenSecret), 10)
	testTokenStruct = &database.Token{
		OwnerId:      testUser.ID,
		TokenId:      tokenId,
		TokenHash:    string(tokenHash),
		Expires:      time.Now().Add(time.Hour * 24),
		TokenPurpose: database.TokenPurposeSession,
		Used:         true,
		Recalled:     false,
	}

	context = NewTestContextBuilder(w).
		SetDefaultUrl().
		SetEndpoint("/auth/sessions").
		SetUser(testUser).
		SetMethod("DELETE").
		SetToken(testToken).
		SetHeader("Accept", "application/json").
		Context

	return w, context, testUser, tokenId, tokenSecret, testTokenStruct
}

func TestAuthHandlersDeleteSessionOk(t *testing.T) {
	w, context, testUser, tokenId, tokenSecret, testTokenStruct := SetupAuthHandlersDeleteSession()

	respBodyExpected := api.DefaultResponse{Status: api.OK}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		Logout(
			gomock.Eq(tokenId),
			gomock.Eq(tokenSecret),
		).
		Return(
			testUser,
			testTokenStruct,
			nil,
		)

	handler.postSession(context)

	respBody, respCode, respParseErr := ExtractResponse[api.PostAccountSessionResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(200), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
	// TODO: MatchEntities and gomock.Eq
}

func TestAuthHandlersDeleteSessionNok_InvTok(t *testing.T) {
	w, context, _, tokenId, tokenSecret, _ := SetupAuthHandlersDeleteSession()

	respBodyExpected := api.DefaultResponse{Status: api.UNAUTHORIZED}

	// test env prep
	ctrl := gomock.NewController(t)
	handler := GetAuthHandlers(ctrl)

	handler.authManager.(*MockAuthManager).
		EXPECT().
		Logout(
			gomock.Eq(tokenId),
			gomock.Eq(tokenSecret),
		).
		Return(
			nil,
			nil,
			managers.InvalidToken,
		)

	handler.postSession(context)

	respBody, respCode, respParseErr := ExtractResponse[api.PostAccountSessionResponse](w)

	require.Nilf(t, respParseErr, "Failed to parse JSON response")
	require.Equalf(t, respCode, int(401), "Response returned unexpected status code")
	require.Truef(t, MatchEntities(respBodyExpected, respBody), "Response returned unexpected body contents")
	// TODO: MatchEntities and gomock.Eq
}
