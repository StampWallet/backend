package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/lithammer/shortuuid/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/utils"
)

var ErrUnknownToken = errors.New("Invalid token")
var ErrTokenExpired = errors.New("Token expired")
var ErrTokenUsed = errors.New("Token used")

// A TokenService is a service for managing tokens. Tokens are used for authorization instead of actual
// user credentials in scenarios where temporary, disposable credentials are desirable.
// Example scenarios: identifying user session after login, identifying the user from a verification email.
type TokenService interface {
	// Creates a new token. Returns database.Token and token secret (hashed secret is stored in the database).
	// Token secret is confidential and should not be stored on the backend.
	// purpose controls Check behavior.
	// If TokenPurpose is TokenPurposeEmail, token is invalidated after Check is called on the token.
	// If TokenPurpose is TokenPurposeSession, token expiration date is changed on each Check call
	// (the date is moved exactly a week from call date, although that could change any time).
	Create(user *User, purpose TokenPurposeEnum, expiration time.Time) (*Token, string, error)

	// Checks if token with tokenId exists in the database.
	// Returns it if tokenSecret matches database.Token.TokenHash.
	Check(tokenId string, tokenSecret string) (*Token, error)

	// Invalidates the token - the token cannot be used after that, Check will return ErrUnknownToken.
	Invalidate(token *Token) (*Token, error)
}

type TokenServiceImpl struct {
	baseServices BaseServices
}

// TODO maybe TokenService should decide about expiration date instead of the caller
// especially since it will handle
func CreateTokenServiceImpl(baseServices BaseServices) *TokenServiceImpl {
	return &TokenServiceImpl{
		baseServices: baseServices,
	}
}

func (service *TokenServiceImpl) Create(user *User, purpose TokenPurposeEnum, expiration time.Time) (*Token, string, error) {
	// Create token secret and hash it
	secret := shortuuid.New()
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), 10)
	if err != nil {
		return nil, "", fmt.Errorf("%s bcrypt failed to generate password: %+v", CallerFilename(), err)
	}

	// Create the token and store it in the database
	token := Token{
		TokenId:      shortuuid.New(),
		TokenHash:    string(hash),
		TokenPurpose: TokenPurposeSession,
		Expires:      expiration,
		Used:         false,
		Recalled:     false,
		User:         user,
	}

	tx := service.baseServices.Database.Create(&token)
	if err := tx.GetError(); err != nil {
		return nil, "", fmt.Errorf("%s database failed to create token: %+v", CallerFilename(), err)
	}

	return &token, secret, nil
}

func (service *TokenServiceImpl) Check(tokenId string, tokenSecret string) (*Token, error) {
	var token Token

	// Get token from the database
	tx := service.baseServices.Database.
		Preload("User").
		First(&token, Token{TokenId: tokenId, Recalled: false})
	err := tx.GetError()
	if err == gorm.ErrRecordNotFound {
		return nil, ErrUnknownToken
	} else if err != nil {
		return nil, fmt.Errorf("%s database failed to find token: %+v", CallerFilename(), err)
	}

	// Check if token is still valid
	if time.Now().After(token.Expires) {
		return nil, ErrTokenExpired
	} else if token.Recalled {
		return nil, ErrUnknownToken
	}

	// Check if TokenHash matches tokenSecret
	err = bcrypt.CompareHashAndPassword([]byte(token.TokenHash), []byte(tokenSecret))
	if err != nil {
		return nil, ErrUnknownToken
	}

	// Check if token is valid for it's purpose
	if token.TokenPurpose == TokenPurposeEmail && token.Used {
		return nil, ErrTokenUsed
	} else if token.TokenPurpose == TokenPurposeSession {
		token.Expires = time.Now().Add(7 * 24 * time.Hour)
	}

	// Update the token in the database
	token.Used = true
	tx = service.baseServices.Database.Save(&token)
	if err := tx.GetError(); err != nil {
		return nil, fmt.Errorf("%s database failed to update token: %+v", CallerFilename(), err)
	}

	return &token, nil
}

func (service *TokenServiceImpl) Invalidate(token *Token) (*Token, error) {
	// Update token in the database
	token.Recalled = true
	tx := service.baseServices.Database.Save(token)
	if err := tx.GetError(); err != nil {
		return nil, fmt.Errorf("%s database failed to update token: %+v", CallerFilename(), err)
	}
	return token, nil
}
