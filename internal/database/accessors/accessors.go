package database

// TODO change to accessors
// package database in package database...

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
		return ErrNotFound
	} else if err != nil {
		return fmt.Errorf("%w: %w", ErrDatabaseError, err)
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
		return empty, ErrNoAccess
	} else {
		return empty, fmt.Errorf("%w: %w", ErrDatabaseError, err)
	}
}

// BusinessAuthorizedAccessor

type BusinessAuthorizedAccessor interface {
	Get(business *Business, cond BusinessOwnedEntity) (BusinessOwnedEntity, error)
	GetAll(business *Business, cond BusinessOwnedEntity) ([]BusinessOwnedEntity, error)
}

type BusinessAuthorizedAccessorImpl struct {
	database GormDB
}

func CreateBusinessAuthorizedAccessorImpl(database GormDB) *BusinessAuthorizedAccessorImpl {
	return &BusinessAuthorizedAccessorImpl{
		database: database,
	}
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

// NOTE shouldnt be used for huge amounts of data
func (accessor *BusinessAuthorizedAccessorImpl) GetAll(business *Business, conds BusinessOwnedEntity) ([]BusinessOwnedEntity, error) {
	condsValue := reflect.ValueOf(conds)
	field := condsValue.Elem().FieldByName("BusinessId")
	if field.IsValid() {
		if field.CanSet() && field.Kind() == reflect.Uint {
			field.SetUint(uint64(business.ID))
		}
	}

	dbResult := reflect.New(reflect.SliceOf(reflect.TypeOf(conds).Elem()))

	tx := accessor.database.Find(dbResult.Interface(), condsValue.Interface())
	if err := tx.GetError(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []BusinessOwnedEntity{}, nil
		}
		return nil, err
	}

	var result []BusinessOwnedEntity
	for i := 0; i != dbResult.Elem().Len(); i += 1 {
		result = append(result, dbResult.Elem().Index(i).Addr().Interface().(BusinessOwnedEntity))
	}

	return result, nil
}

// UserAuthorizedAccessor

type UserAuthorizedAccessor interface {
	Get(user *User, cond UserOwnedEntity) (UserOwnedEntity, error)
	GetAll(user *User, cond UserOwnedEntity, preloads []string) ([]UserOwnedEntity, error)
}

type UserAuthorizedAccessorImpl struct {
	database GormDB
}

func CreateUserAuthorizedAccessorImpl(database GormDB) *UserAuthorizedAccessorImpl {
	return &UserAuthorizedAccessorImpl{
		database: database,
	}
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

// NOTE shouldnt be used for huge amounts of data
func (accessor *UserAuthorizedAccessorImpl) GetAll(user *User, conds UserOwnedEntity, preloads []string) ([]UserOwnedEntity, error) {
	condsValue := reflect.ValueOf(conds)
	field := condsValue.Elem().FieldByName("OwnerId")
	if field.IsValid() {
		if field.CanSet() && field.Kind() == reflect.Uint {
			field.SetUint(uint64(user.ID))
		}
	}

	dbResult := reflect.New(reflect.SliceOf(reflect.TypeOf(conds).Elem()))
	tx := accessor.database
	for _, v := range preloads {
		tx = tx.Preload(v)
	}
	tx = tx.Find(dbResult.Interface(), condsValue.Interface())
	if err := tx.GetError(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []UserOwnedEntity{}, nil
		}
		return nil, err
	}

	var result []UserOwnedEntity
	for i := 0; i != dbResult.Elem().Len(); i += 1 {
		result = append(result, dbResult.Elem().Index(i).Addr().Interface().(UserOwnedEntity))
	}

	return result, nil
}

// AuthorizedTransactionAccessor

type AuthorizedTransactionAccessor interface {
	GetForBusiness(business *Business, transactionCode string) (*Transaction, error)
	GetForUser(user *User, transactionCode string) (*Transaction, error)
}

type AuthorizedTransactionAccessorImpl struct {
	database GormDB
}

func CreateAuthorizedTransactionAccessorImpl(database GormDB) *AuthorizedTransactionAccessorImpl {
	return &AuthorizedTransactionAccessorImpl{
		database: database,
	}
}

func (accessor *AuthorizedTransactionAccessorImpl) GetForBusiness(business *Business, transactionCode string) (*Transaction, error) {
	var transaction Transaction
	tx := accessor.database.
		//NOTE I'm not very confident about the efficiency of these preloads
		// On the other hand, this accessor is currently used in few specific situations
		// where these preloads are actually useful (at least some fields...
		// So maybe it's not a good idea to optimize now.
		// The dataset is never very big and rows by themselves aren't very big either.
		// Still, double join
		Preload("TransactionDetails").
		Preload("TransactionDetails.OwnedItem").
		Preload("TransactionDetails.ItemDefinition").
		First(&transaction, Transaction{
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
