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
			Database: GetTestDatabase(),
		},
	}
}

// Tests

// Test TokenServiceImpl.Create on happy path
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
		"Token expiration date be at least 23 hours after now")
	require.Truef(t, time.Now().Add(25*time.Hour).After(token.Expires),
		"Token expiration date be should be at most 25 hours after now")

	// Check if token was propoery saved into the database
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
		"Token expiration date be at least 23 hours after now")
	require.Truef(t, time.Now().Add(25*time.Hour).After(dbToken.Expires),
		"Token expiration date be should be at most 25 hours after now")
}

// Test TokenServiceImpl.Check with TokenPurposeEmail on happy path and after the token was already used
func TestTokenServiceCheckValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := GetTokenService(ctrl)
	user := GetTestUser(service.baseServices.Database)
	token, secret := GetTestToken(service.baseServices.Database, user) //returns email token

	nToken, err := service.Check(token.TokenId, secret)
	require.Nilf(t, err, "TokenService.Check should not return an error")
	require.NotNilf(t, nToken.User, "User should not be nil")
	require.NotNilf(t, nToken, "Token should not be nil")
	require.Equalf(t, token.TokenId, nToken.TokenId, "TokenService.Check should return the expected token")
	require.Equalf(t, true, nToken.Used, "Token should be marked as used")
	require.Equalf(t, nToken.User.PublicId, user.PublicId, "TokenService.Check should return the expected user")

	_, err = service.Check(token.TokenId, secret)
	require.Equalf(t, err, ErrTokenUsed, "TokenService.Check should return ErrTokenUsed on used email token")
}

// Test TokenServiceImpl.Check with TokenPurposeSession on happy path
func TestTokenServiceCheckValidSessionToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := GetTokenService(ctrl)
	user := GetTestUser(service.baseServices.Database)
	token, secret := GetTestSessionToken(service.baseServices.Database, user, time.Now().Add(time.Hour))

	nToken, err := service.Check(token.TokenId, secret)
	require.Nilf(t, err, "TokenService.Check should not return an error")
	require.NotNilf(t, nToken.User, "User should not be nil")
	require.NotNilf(t, nToken, "Token should not be nil")
	require.Equalf(t, token.TokenId, nToken.TokenId, "TokenService.Check should return the expected token")
	require.Equalf(t, true, nToken.Used, "Token should be marked as used")
	require.Equalf(t, nToken.User.PublicId, user.PublicId, "TokenService.Check should return the expected user")
	require.Truef(t, time.Now().Add((7*24*time.Hour)-time.Hour).Before(nToken.Expires), "session token expiration date should be updated")

	_, err = service.Check(token.TokenId, secret)
	require.Nilf(t, err, "TokenService.Check should not return ErrTokenUsed on used session token")
}

// Test TokenServiceImpl.Check when tokenId is invalid
func TestTokenServiceCheckInvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := GetTokenService(ctrl)
	user := GetTestUser(service.baseServices.Database)
	_, secret := GetTestToken(service.baseServices.Database, user)

	nToken, err := service.Check("invalid id", secret)
	require.Equalf(t, ErrUnknownToken, err, "TokenService.Check should return a UnknownToken error")
	require.Nilf(t, nToken, "TokenService.Check should return a nil token")
}

// Test TokenServiceImpl.Check when tokenId is valid, but token secret does not match
func TestTokenServiceCheckInvalidSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := GetTokenService(ctrl)
	user := GetTestUser(service.baseServices.Database)
	token, secret := GetTestToken(service.baseServices.Database, user)

	nToken, err := service.Check(token.TokenId, secret+"invalid")
	require.Equalf(t, ErrUnknownToken, err, "TokenService.Check should return a UnknownToken error")
	require.Nilf(t, nToken, "TokenService.Check should return a nil token")
}

// Test TokenServiceImpl.Invalidate on happy path
func TestTokenServiceInvalidate(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := GetTokenService(ctrl)
	user := GetTestUser(service.baseServices.Database)
	token, secret := GetTestToken(service.baseServices.Database, user)

	nToken, err := service.Invalidate(token)
	require.Nilf(t, err, "TokenService.Invalidate should return nil error")
	require.NotNilf(t, nToken.User, "User should not be nil")
	require.NotNilf(t, nToken, "Token should not be nil")
	require.Equalf(t, token.TokenId, nToken.TokenId, "TokenService.Invalidate should return the expected token")
	require.Equalf(t, true, nToken.Recalled, "TokenService.Invalidate should return a recalled token")
	require.Equalf(t, user.PublicId, nToken.User.PublicId, "TokenService.Invalidate should return the expected user")

	testToken, err := service.Check(token.TokenId, secret)
	require.Equalf(t, ErrUnknownToken, err, "TokenService.Check should return an error")
	require.Nilf(t, testToken, "TokenService.Check should return a nil token")

	var dbToken Token
	tx := service.baseServices.Database.First(&dbToken, Token{TokenId: token.TokenId})
	require.Nilf(t, tx.GetError(), "Database.First should return a nil error")
	require.Equalf(t, true, dbToken.Recalled, "Token in the database should be recalled")
}
