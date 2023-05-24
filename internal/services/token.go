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

type TokenService interface {
	Create(user *User, purpose TokenPurposeEnum, expiration time.Time) (*Token, string, error)
	Check(tokenId string, tokenSecret string) (*Token, error)
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
	secret := shortuuid.New()
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), 10)
	if err != nil {
		return nil, "", fmt.Errorf("%s bcrypt failed to generate password: %+v", CallerFilename(), err)
	}
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
	tx := service.baseServices.Database.
		Preload("User").
		First(&token, Token{TokenId: tokenId, Recalled: false})
	err := tx.GetError()
	if err == gorm.ErrRecordNotFound {
		return nil, ErrUnknownToken
	} else if err != nil {
		return nil, fmt.Errorf("%s database failed to find token: %+v", CallerFilename(), err)
	}

	if time.Now().After(token.Expires) {
		return nil, ErrTokenExpired
	}

	if token.Recalled {
		return nil, ErrUnknownToken
	}

	err = bcrypt.CompareHashAndPassword([]byte(token.TokenHash), []byte(tokenSecret))
	if err != nil {
		return nil, ErrUnknownToken
	}

	if token.TokenPurpose == TokenPurposeEmail && token.Used {
		return nil, ErrTokenUsed
	} else if token.TokenPurpose == TokenPurposeSession {
		token.Expires = time.Now().Add(7 * 24 * time.Hour)
	}

	token.Used = true
	tx = service.baseServices.Database.Save(&token)
	if err := tx.GetError(); err != nil {
		return nil, fmt.Errorf("%s database failed to update token: %+v", CallerFilename(), err)
	}

	return &token, nil
}

func (service *TokenServiceImpl) Invalidate(token *Token) (*Token, error) {
	token.Recalled = true
	tx := service.baseServices.Database.Save(token)
	if err := tx.GetError(); err != nil {
		return nil, fmt.Errorf("%s database failed to update token: %+v", CallerFilename(), err)
	}
	return token, nil
}
