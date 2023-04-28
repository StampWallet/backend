package database

import (
	"errors"
	"reflect"

	. "github.com/StampWallet/backend/internal/database"
	"gorm.io/gorm"
)

var NotFound = errors.New("Entity not found")
var NoAccess = errors.New("No access to entity")

type BusinessAuthorizedAccessor interface {
	Get(business *Business, cond BusinessOwnedEntity) (BusinessOwnedEntity, error)
}

type BusinessAuthorizedAccessorImpl struct {
	database GormDB
}

func (accessor *BusinessAuthorizedAccessorImpl) Get(business *Business, conds BusinessOwnedEntity) (BusinessOwnedEntity, error) {
	result := reflect.New(reflect.TypeOf(conds).Elem()).Interface().(BusinessOwnedEntity)
	tx := accessor.database.First(&result, conds)
	err := tx.GetError()
	if err != nil {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, NotFound
	}
	if tx.GetRowsAffected() != 1 {
		return nil, NotFound
	}
	if result.GetBusinessId() == business.ID {
		return result, nil
	} else {
		return nil, NoAccess
	}
}

type UserAuthorizedAccessor interface {
	Get(user *User, cond UserOwnedEntity) (UserOwnedEntity, error)
}

type UserAuthorizedAccessorImpl struct {
	database GormDB
}

func (accessor *UserAuthorizedAccessorImpl) Get(user *User, cond UserOwnedEntity) (UserOwnedEntity, error) {
	return nil, nil
}

type AuthorizedTransactionAccessor interface {
	GetForBusiness(business *Business, transactionCode string) (*Transaction, error)
	GetForUser(user *User, transactionCode string) (*Transaction, error)
}

type AuthorizedTransactionAccessorImpl struct {
	database GormDB
}

// NOTE remember to load all transaction details
func (accessor *AuthorizedTransactionAccessorImpl) GetForBusiness(business *Business, transactionCode string) (*Transaction, error) {
	return nil, nil
}

func (accessor *AuthorizedTransactionAccessorImpl) GetForUser(user *User, transactionCode string) (*Transaction, error) {
	return nil, nil
}

type UserOwnedEntity interface {
	GetUserId() uint
}

type BusinessOwnedEntity interface {
	GetBusinessId() uint
}
