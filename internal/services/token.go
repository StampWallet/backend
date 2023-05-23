package services

import (
	"errors"
	"time"

	. "github.com/StampWallet/backend/internal/database"
)

var UnknownToken = errors.New("Invalid token")

type TokenService interface {
	Create(user *User, purpose TokenPurposeEnum, expiration time.Time) (*Token, string, error)
	Check(tokenId string, tokenSecret string) (*User, *Token, error)
	Invalidate(token *Token) (*User, *Token, error)
}

type TokenServiceImpl struct {
	baseServices BaseServices
}

func CreateTokenServiceImpl(baseServices BaseServices) *TokenServiceImpl {
	return &TokenServiceImpl{
		baseServices: baseServices,
	}
}

func (service *TokenServiceImpl) Create(user *User, purpose TokenPurposeEnum, expiration time.Time) (*Token, string, error) {
	return nil, "", nil
}

func (service *TokenServiceImpl) Check(tokenId string, tokenSecret string) (*User, *Token, error) {
	return nil, nil, nil
}

func (service *TokenServiceImpl) Invalidate(token *Token) (*User, *Token, error) {
	return nil, nil, nil
}
