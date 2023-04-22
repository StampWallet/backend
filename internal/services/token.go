package services

import (
    "time"
    . "github.com/StampWallet/backend/internal/database"
)

type TokenService interface {
    Create(user User, purpose TokenPurposeEnum, expiration time.Time) (Token, error)
    Check(tokenId string, token string) (User, error)
    Invalidate(tokenId string) (User, error)
}

type TokenServiceImpl struct {
    baseServices BaseServices
}
