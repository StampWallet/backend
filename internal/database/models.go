package database

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type ActionTypeEnum string

//TODO inconsistent naming. recalled vs withdrawn

const (
	NoActionType        ActionTypeEnum = "NO_ACTION"
	RedeemedActionType                 = "REDEEMED"
	RecalledActionType                 = "RECALLED"
	CancelledActionType                = "CANCELLED"
)

type TransactionStateEnum string

const (
	TransactionStateStarted   TransactionStateEnum = "STARTED"
	TransactionStateProcesing                      = "PROCESSING"
	TransactionStateFinished                       = "FINISHED"
	TransactionStateExpired                        = "EXPIRED"
)

type TokenPurposeEnum string

const (
	TokenPurposeSession TokenPurposeEnum = "SESSION"
	TokenPurposeEmail   TokenPurposeEnum = "EMAIL"
)

type OwnedItemStatusEnum string

const (
	OwnedItemStatusOwned     OwnedItemStatusEnum = "OWNED"
	OwnedItemStatusUsed                          = "USED"
	OwnedItemStatusWithdrawn                     = "WITHDRAWN"
	OwnedItemStatusReturned                      = "RETURNED"
)

// MODELS

// LocalCard

type LocalCard struct {
	gorm.Model
	PublicId string
	OwnerId  uint
	Type     string
	Code     string
	Name     string
}

func (entity *LocalCard) GetUserId(_ GormDB) (uint, error) {
	return entity.OwnerId, nil
}

// Token

type Token struct {
	gorm.Model
	OwnerId      uint
	TokenId      string
	TokenHash    string
	Expires      time.Time
	TokenPurpose TokenPurposeEnum
	Used         bool
	Recalled     bool
}

func (entity *Token) GetUserId(_ GormDB) (uint, error) {
	return entity.OwnerId, nil
}

// FileMetadata

type FileMetadata struct {
	gorm.Model
	PublicId string
	OwnerId  uint
	Uploaded sql.NullTime
}

func (entity *FileMetadata) GetUserId(_ GormDB) (uint, error) {
	return entity.OwnerId, nil
}

// User

type User struct {
	gorm.Model
	PublicId      string
	FirstName     string
	LastName      string
	Email         string
	PasswordHash  string
	EmailVerified bool

	Tokens        []Token        `gorm:"foreignkey:OwnerId"`
	LocalCards    []LocalCard    `gorm:"foreignkey:OwnerId"`
	VirtualCards  []VirtualCard  `gorm:"foreignkey:OwnerId"`
	FilesMetadata []FileMetadata `gorm:"foreignkey:OwnerId"`
	Business      Business       `gorm:"foreignkey:OwnerId"`
}

// Business

type Business struct {
	gorm.Model
	PublicId       string
	OwnerId        uint
	Name           string
	Description    string
	Address        string
	GPSCoordinates string
	NIP            string
	KRS            string
	REGON          string
	OwnerName      string
	BannerImageId  string
	IconImageId    string

	ItemDefinitions []ItemDefinition
	MenuItems       []MenuItem
	VirtualCards    []VirtualCard
}

func (entity *Business) GetUserId() (uint, error) {
	return entity.OwnerId, nil
}

// ItemDefinition

type ItemDefinition struct {
	gorm.Model
	PublicId    string
	BusinessId  uint
	Name        string
	Price       uint
	Description string
	ImageId     string
	StartDate   time.Time
	EndDate     time.Time
	MaxAmount   uint
	Available   bool
	Withdrawn   bool
}

func (entity *ItemDefinition) GetBusinessId(_ GormDB) (uint, error) {
	return entity.BusinessId, nil
}

// MenuItem

type MenuItem struct {
	gorm.Model
	BusinessId uint
	FileId     uint
}

func (entity *MenuItem) GetBusinessId() (uint, error) {
	return entity.BusinessId, nil
}

// OwnedItem

type OwnedItem struct {
	gorm.Model
	PublicId      string
	DefinitionId  uint
	VirtualCardId uint
	Used          sql.NullTime
	Status        OwnedItemStatusEnum
}

func (entity *OwnedItem) GetUserId(db GormDB) (uint, error) {
	var virtualCard VirtualCard
	tx := db.First(&virtualCard, VirtualCard{Model: gorm.Model{ID: entity.VirtualCardId}})
	if err := tx.GetError(); err != nil {
		return 0, err
	}
	return virtualCard.OwnerId, nil
}

func (entity *OwnedItem) GetBusinessId(db GormDB) (uint, error) {
	var virtualCard VirtualCard
	tx := db.First(&virtualCard, VirtualCard{Model: gorm.Model{ID: entity.VirtualCardId}})
	if err := tx.GetError(); err != nil {
		return 0, err
	}
	return virtualCard.BusinessId, nil
}

// VirtualCard

type VirtualCard struct {
	gorm.Model
	OwnerId    uint
	PublicId   string
	BusinessId uint
	Points     uint

	Transactions []Transaction
	OwnedItems   []OwnedItem
	Business     Business `gorm:"foreignkey:BusinessId"`
	User         User     `gorm:"foreignkey:OwnerId"`
}

func (entity *VirtualCard) GetUserId(_ GormDB) (uint, error) {
	return entity.OwnerId, nil
}

func (entity *VirtualCard) GetBusinessId(_ GormDB) (uint, error) {
	return entity.BusinessId, nil
}

// Transaction

type Transaction struct {
	gorm.Model
	PublicId      string
	VirtualCardId uint   `gorm:"index:code,unique"`
	Code          string `gorm:"index:code,unique"`
	State         TransactionStateEnum
	AddedPoints   uint

	TransactionDetails []TransactionDetail
	VirtualCard        VirtualCard `gorm:"foreignkey:VirtualCardId"`
}

func (entity *Transaction) GetUserId(db GormDB) (uint, error) {
	var virtualCard VirtualCard
	tx := db.First(&virtualCard, VirtualCard{Model: gorm.Model{ID: entity.VirtualCardId}})
	if err := tx.GetError(); err != nil {
		return 0, err
	}
	return virtualCard.OwnerId, nil
}

func (entity *Transaction) GetBusinessId(db GormDB) (uint, error) {
	var virtualCard VirtualCard
	tx := db.First(&virtualCard, VirtualCard{Model: gorm.Model{ID: entity.VirtualCardId}})
	if err := tx.GetError(); err != nil {
		return 0, err
	}
	return virtualCard.BusinessId, nil
}

// TransactionDetail

type TransactionDetail struct {
	gorm.Model
	TransactionId uint
	ItemId        uint
	Action        ActionTypeEnum
}
