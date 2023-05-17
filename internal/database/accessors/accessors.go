package database

import (
	"errors"
	"fmt"
	"reflect"

	. "github.com/StampWallet/backend/internal/database"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("Entity not found")
var ErrNoAccess = errors.New("No access to entity")
var ErrDatabaseError = errors.New("Database error")

// rule of thumb with managers and accessors: if a manager method requires an object, this probably means, that the object has to be retrived with an accessor. and the accessor will check the (simple, in case of this app) permissions
// if a manager requires only object id, that means that the action can be done by anyone
// handlers should not access the database directly. thats what managers and accessors are for

func checkErr(tx GormDB) error {
	err := tx.GetError()
	if err == gorm.ErrRecordNotFound || tx.GetRowsAffected() != 1 {
		return NotFound
	} else if err != nil {
		return fmt.Errorf("%w: %w", DatabaseError, err)
	}
	return nil
}

type authModel interface {
	OwnedEntity | *Transaction
}

func checkEq[T authModel](el T, expectedId uint, id uint, err error) (T, error) {
	var empty T
	if expectedId == id {
		return el, nil
	} else if expectedId != id {
		return empty, NoAccess
	} else {
		return empty, fmt.Errorf("%w: %w", DatabaseError, err)
	}
}

// BusinessAuthorizedAccessor

type BusinessAuthorizedAccessor interface {
	Get(business *Business, cond BusinessOwnedEntity) (BusinessOwnedEntity, error)
}

type BusinessAuthorizedAccessorImpl struct {
	database GormDB
}

func (accessor *BusinessAuthorizedAccessorImpl) Get(business *Business, conds BusinessOwnedEntity) (BusinessOwnedEntity, error) {
	result := reflect.New(reflect.TypeOf(conds).Elem()).Interface().(BusinessOwnedEntity)
	tx := accessor.database.First(&result, conds)
	if err := checkErr(tx); err != nil {
		return nil, err
	}

	id, err := result.GetBusinessId(accessor.database)
	return checkEq(result, business.ID, id, err)
}

// UserAuthorizedAccessor

type UserAuthorizedAccessor interface {
	Get(user *User, cond UserOwnedEntity) (UserOwnedEntity, error)
}

type UserAuthorizedAccessorImpl struct {
	database GormDB
}

func (accessor *UserAuthorizedAccessorImpl) Get(user *User, conds UserOwnedEntity) (UserOwnedEntity, error) {
	result := reflect.New(reflect.TypeOf(conds).Elem()).Interface().(UserOwnedEntity)
	tx := accessor.database.First(&result, conds)
	if err := checkErr(tx); err != nil {
		return nil, err
	}

	id, err := result.GetUserId(accessor.database)
	return checkEq(result, user.ID, id, err)
}

// AuthorizedTransactionAccessor

type AuthorizedTransactionAccessor interface {
	GetForBusiness(business *Business, transactionCode string) (*Transaction, error)
	GetForUser(user *User, transactionCode string) (*Transaction, error)
}

type AuthorizedTransactionAccessorImpl struct {
	database GormDB
}

func (accessor *AuthorizedTransactionAccessorImpl) GetForBusiness(business *Business, transactionCode string) (*Transaction, error) {
	var transaction Transaction
	tx := accessor.database.Preload("TransactionDetails").First(&transaction, Transaction{
		Code:        transactionCode,
		VirtualCard: &VirtualCard{Business: business},
	})
	if err := checkErr(tx); err != nil {
		return nil, err
	}

	id, err := transaction.GetBusinessId(accessor.database)
	return checkEq(&transaction, business.ID, id, err)
}

func (accessor *AuthorizedTransactionAccessorImpl) GetForUser(user *User, transactionCode string) (*Transaction, error) {
	var transaction Transaction
	tx := accessor.database.Preload("TransactionDetails").First(&transaction, Transaction{
		Code:        transactionCode,
		VirtualCard: &VirtualCard{User: user},
	})
	if err := checkErr(tx); err != nil {
		return nil, err
	}

	id, err := transaction.GetUserId(accessor.database)
	return checkEq(&transaction, user.ID, id, err)
}

// OwnedEntities

type OwnedEntity interface {
}

type UserOwnedEntity interface {
	OwnedEntity
	GetUserId(db GormDB) (uint, error)
}

type BusinessOwnedEntity interface {
	OwnedEntity
	GetBusinessId(db GormDB) (uint, error)
}
