package database

import (
	"time"
	"database/sql"
	"gorm.io/gorm"
)

type ActionTypeEnum string

const (
	RedeemedActionType ActionTypeEnum	= "REDEEMED"
	RecalledActionType					= "RECALLED"
	CancelledActionType					= "CANCELLED"
)

type TransactionStateEnum string

const (
	StartedTransactionState TransactionStateEnum	= "STARTED"
	ProcessedTransactionState						= "PROCESSED"
	FinishedTransactionState						= "FINISHED"
	ExpiredTransactionState							= "EXPIRED"
)

type TokenPurposeEnum string

const (
	SessionTokenPurpose TokenPurposeEnum	= "SESSION"
	EmailTokenPurpose TokenPurposeEnum		= "EMAIL"
)

type OwnedItemStatusEnum string

const (
	OwnedItemStatus OwnedItemStatusEnum	= "OWNED"
	UsedItemStatus						= "USED"
	RecalledItemStatus					= "RECALLED"
)


type LocalCard struct {
	gorm.Model
	PublicId string
	OwnerId uint
	Type string
	Code string
	Name string
}

func (entity *LocalCard) GetUserId() uint {
	return 0
}

type Token struct {
	gorm.Model
	OwnerId uint
	TokenId string
	TokenHash string
	Expires time.Time
	TokenPurpose TokenPurposeEnum
	Used bool
	Recalled bool
}

func (entity *Token) GetUserId() uint {
	return 0
}

type FileMetadata struct {
	gorm.Model
	PublicId string
	OwnerId uint
	Uploaded sql.NullTime
}

func (entity *FileMetadata) GetUserId() uint {
	return 0
}

type User struct {
	gorm.Model
	PublicId string
	FirstName string
	LastName string
	Email string
	PasswordHash string
	EmailVerified bool

	Tokens []Token `gorm:"foreignkey:OwnerId"`
	LocalCards []LocalCard `gorm:"foreignkey:OwnerId"`
	VirtualCards []VirtualCard `gorm:"foreignkey:OwnerId"`
	FilesMetadata []FileMetadata `gorm:"foreignkey:OwnerId"`
	Business Business `gorm:"foreignkey:OwnerId"`
}

type Business struct {
	gorm.Model
	PublicId string
	OwnerId uint
	Name string
	Description string
	Address string
	GPSCoordinates string
	NIP string
	KRS string
	REGON string
	OwnerName string
	BannerImageId string
	IconImageId string

	ItemDefinitions []ItemDefinition
	MenuItems []MenuItem
	VirtualCards []VirtualCard
}

func (entity *Business) GetUserId() uint {
	return 0
}

type ItemDefinition struct {
	gorm.Model
	PublicId string
	BusinessId uint
	Name string
	Price uint
	Description string
	ImageId string
	StartDate time.Time
	EndDate time.Time
	MaxAmount uint
	Available bool
}

func (entity *Business) GetBusinessId() uint {
	return 0
}

type MenuItem struct {
	gorm.Model
	BusinessId uint
	FileId uint
}

func (entity *MenuItem) GetBusinessId() uint {
	return 0
}

type OwnedItem struct {
	gorm.Model
	PublicId string
	DefinitionId uint
	VirtualCardId uint
	Used sql.NullTime
	Status OwnedItemStatusEnum
}

func (entity *OwnedItem) GetUserId() uint {
	return 0
}

type VirtualCard struct {
	gorm.Model
	OwnerId uint
	PublicId string
	BusinessId uint
	Points uint

	Transactions []Transaction
	OwnedItems []OwnedItem
}

type Transaction struct {
	gorm.Model
	PublicId string
	VirtualCardId uint
	Code string
	State TransactionStateEnum
	AddedPoints uint

	TransactionsDetails []TransactionDetails
}

type TransactionDetails struct {
	gorm.Model
	TransactionId uint
	ItemId uint
	Action ActionTypeEnum
}
