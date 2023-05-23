package services

import (
	"log"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/testutils"
)

func GetTokenService(ctrl *gomock.Controller) *TokenServiceImpl {
	return &TokenServiceImpl{
		baseServices: BaseServices{
			Logger:   log.Default(),
			Database: GetDatabase(),
		},
	}
}

// Tests

func TestTokenServiceCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := GetTokenService(ctrl)
	user := GetTestUser(service.baseServices.Database)

	token, secret, err := service.Create(user, TokenPurposeSession, time.Now().Add(24*time.Hour))
	require.Nilf(t, err, "Error should be nil")
	require.NotNilf(t, token, "Token should not be nil")
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(token.TokenHash), []byte(secret))
	require.Nilf(t, bcryptErr, "Bcrypt validation should pass")
	require.Equalf(t, user.ID, token.OwnerId, "Token ownerId should match")
	require.Equalf(t, token.TokenPurpose, TokenPurposeSession, "Token purpose should be session")
	require.Falsef(t, token.Used, "Token used should be false")
	require.Falsef(t, token.Recalled, "Token recalled should be false")
	require.Truef(t, time.Now().Before(token.Expires.Add(-23*time.Hour)),
		"Token expiration date should way before now")

	var dbToken Token
	tx := service.baseServices.Database.First(&dbToken, Token{TokenId: token.TokenId})
	require.Nilf(t, tx.GetError(), "Database.First should not return an error")

	bcryptErr = bcrypt.CompareHashAndPassword([]byte(dbToken.TokenHash), []byte(secret))
	require.Nilf(t, bcryptErr, "Bcrypt validation should pass")
	require.Equalf(t, user.ID, dbToken.OwnerId, "Token ownerId should match")
	require.Falsef(t, dbToken.Used, "Token used should be false")
	require.Falsef(t, dbToken.Recalled, "Token recalled should be false")
	require.Equalf(t, dbToken.TokenPurpose, TokenPurposeSession, "Token purpose should be session")
	require.Truef(t, time.Now().Before(dbToken.Expires.Add(-23*time.Hour)),
		"Token expiration date should way before now")
}

func TestTokenServiceCheckValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := GetTokenService(ctrl)
	user := GetTestUser(service.baseServices.Database)
	token, secret := GetTestToken(service.baseServices.Database, user) //returns email token

	nUser, nToken, err := service.Check(token.TokenId, secret)
	require.Nilf(t, err, "TokenService.Check should not return an error")
	require.NotNilf(t, nUser, "User should not be nil")
	require.NotNilf(t, nToken, "Token should not be nil")
	require.Equalf(t, token.TokenId, nToken.TokenId, "TokenService.Check should return the expected token")
	require.Equalf(t, true, nToken.Used, "Token should be marked as used")

	require.Equalf(t, nUser.PublicId, user.PublicId, "TokenService.Check should return the expected user")
}

func TestTokenServiceCheckInvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := GetTokenService(ctrl)
	user := GetTestUser(service.baseServices.Database)
	_, secret := GetTestToken(service.baseServices.Database, user)

	nUser, nToken, err := service.Check("invalid id", secret)
	require.ErrorAsf(t, err, UnknownToken, "TokenService.Check should return a UnknownToken error")
	require.Nilf(t, nUser, "TokenService.Check should return a nil user")
	require.Nilf(t, nToken, "TokenService.Check should return a nil token")
}

func TestTokenServiceCheckInvalidSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := GetTokenService(ctrl)
	user := GetTestUser(service.baseServices.Database)
	token, secret := GetTestToken(service.baseServices.Database, user)

	nUser, nToken, err := service.Check(token.TokenId, secret+"invalid")
	require.ErrorAsf(t, err, UnknownToken, "TokenService.Check should return a UnknownToken error")
	require.Nilf(t, nUser, "TokenService.Check should return a nil user")
	require.Nilf(t, nToken, "TokenService.Check should return a nil token")
}

func TestTokenServiceInvalidate(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := GetTokenService(ctrl)
	user := GetTestUser(service.baseServices.Database)
	token, secret := GetTestToken(service.baseServices.Database, user)

	nUser, nToken, err := service.Invalidate(token)
	require.Nilf(t, err, "TokenService.Invalidate should return nil error")
	require.NotNilf(t, nUser, "User should not be nil")
	require.NotNilf(t, nToken, "Token should not be nil")
	require.Equalf(t, token.TokenId, nToken.TokenId, "TokenService.Invalidate should return the expected token")
	require.Equalf(t, true, nToken.Recalled, "TokenService.Invalidate should return a recalled token")
	require.Equalf(t, user.PublicId, nUser.PublicId, "TokenService.Invalidate should return the expected user")

	testUser, testToken, err := service.Check(token.TokenId, secret)
	require.ErrorAsf(t, err, UnknownToken, "TokenService.Check should return an error")
	require.Nilf(t, testUser, "TokenService.Check should return a nil user")
	require.Nilf(t, testToken, "TokenService.Check should return a nil token")

	var dbToken Token
	tx := service.baseServices.Database.First(&dbToken, Token{TokenId: token.TokenId})
	require.Nilf(t, tx.GetError(), "Database.First should return a nil error")
	require.Equalf(t, true, dbToken.Recalled, "Token in the database should be recalled")
}
