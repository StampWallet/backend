package services

import (
    "time"
    "errors"
    . "github.com/StampWallet/backend/internal/database"
)

var UnknownToken = errors.New("Invalid token")

type TokenService interface {
    Create(user User, purpose TokenPurposeEnum, expiration time.Time) (*Token, error)
    Check(tokenId string, token string) (*User, *Token, error)
    Invalidate(token Token) (*User, *Token, error)
}

type TokenServiceImpl struct {
    baseServices BaseServices
}
