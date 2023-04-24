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
	EmailTokenPurpose						= "EMAIL"
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
	OwnerId uint64
	Type string
	Code string
	Name string
}

func (entity *LocalCard) GetUserId() uint64 {
	return 0
}

type Token struct {
	gorm.Model
	OwnerId uint64
	TokenId string
	TokenHash string
	Expires time.Time
	TokenPurpose TokenPurposeEnum
	Used bool
	Recalled bool
}

func (entity *Token) GetUserId() uint64 {
	return 0
}

type FileMetadata struct {
	gorm.Model
	PublicId string
	OwnerId uint64
	Uploaded sql.NullTime
}

func (entity *FileMetadata) GetUserId() uint64 {
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
	OwnerId uint64
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

func (entity *Business) GetUserId() uint64 {
	return 0
}

type ItemDefinition struct {
	gorm.Model
	PublicId string
	BusinessId uint64
	Name string
	Price uint64
	Description string
	ImageId string
	StartDate time.Time
	EndDate time.Time
	MaxAmount uint64
	Available bool
}

func (entity *Business) GetBusinessId() uint64 {
	return 0
}

type MenuItem struct {
	gorm.Model
	BusinessId uint64
	FileId uint64
}

func (entity *MenuItem) GetBusinessId() uint64 {
	return 0
}

type OwnedItem struct {
	gorm.Model
	PublicId string
	DefinitionId uint64
	VirtualCardId uint64
	Used sql.NullTime
	Status OwnedItemStatusEnum
}

func (entity *OwnedItem) GetUserId() uint64 {
	return 0
}

type VirtualCard struct {
	gorm.Model
	OwnerId uint64
	PublicId string
	BusinessId uint64
	Points uint64

	Transactions []Transaction
	OwnedItems []OwnedItem
}

type Transaction struct {
	gorm.Model
	PublicId string
	VirtualCardId uint64
	Code string
	State TransactionStateEnum
	AddedPoints uint64

	TransactionsDetails []TransactionDetails
}

type TransactionDetails struct {
	gorm.Model
	TransactionId uint64
	ItemId uint64
	Action ActionTypeEnum
}
