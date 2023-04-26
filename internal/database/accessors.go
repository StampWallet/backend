package database

import (
	"gorm.io/gorm"
)

type BusinessAuthorizedAccessor[T BusinessOwnedEntity] interface {
	Get(business *Business, publicId string) (*T, error)
}

type BusinessAuthorizedAccessorImpl[T BusinessOwnedEntity] struct {
	database gorm.DB
}

func (accessor *BusinessAuthorizedAccessorImpl[T]) Get(business *Business, publicId string) (*T, error) {
	return nil, nil
}

type UserAuthorizedAccessor[T UserOwnedEntity] interface {
	Get(user *User, publicId string) (T, error)
}

type UserAuthorizedAccessorImpl[T UserOwnedEntity] struct {
	database gorm.DB
}

func (accessor *UserAuthorizedAccessorImpl[T]) Get(user *User, publicId string) (*T, error) {
	return nil, nil
}

type AuthorizedTransactionAccessor interface {
	AccessFromBusiness(user *User, transactionCode string) (*Transaction, error)
	AccessFromUser(user *User, transactionCode string) (*Transaction, error)
}

type AuthorizedTransactionAccessorImpl struct {
	database gorm.DB
}

func (accessor *AuthorizedTransactionAccessorImpl) AccessFromBusiness(user User, transactionCode string) (*Transaction, error) {
	return nil, nil
}

func (accessor *AuthorizedTransactionAccessorImpl) AccessFromUser(user User, transactionCode string) (*Transaction, error) {
	return nil, nil
}

type UserOwnedEntity interface {
	GetUserId() uint64
}

type BusinessOwnedEntity interface {
	GetBusinessId() uint64
}

type TransactionOwnedEntity interface {
	GetTransactionId() uint64
}
