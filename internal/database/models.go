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
)

type LocalCard struct {
	gorm.Model
	PublicId string
	OwnerId  uint
	Type     string
	Code     string
	Name     string
}

func (entity *LocalCard) GetUserId() uint {
	return 0
}

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

func (entity *Token) GetUserId() uint {
	return 0
}

type FileMetadata struct {
	gorm.Model
	PublicId string
	OwnerId  uint
	Uploaded sql.NullTime
}

func (entity *FileMetadata) GetUserId() uint {
	return 0
}

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

func (entity *Business) GetUserId() uint {
	return 0
}

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

func (entity *ItemDefinition) GetBusinessId() uint {
	return entity.BusinessId
}

type MenuItem struct {
	gorm.Model
	BusinessId uint
	FileId     uint
}

func (entity *MenuItem) GetBusinessId() uint {
	return entity.BusinessId
}

type OwnedItem struct {
	gorm.Model
	PublicId      string
	DefinitionId  uint
	VirtualCardId uint
	Used          sql.NullTime
	Status        OwnedItemStatusEnum
}

func (entity *OwnedItem) GetUserId() uint {
	return 0
}

func (entity *OwnedItem) GetBusinessId() uint {
	return 0
}

type VirtualCard struct {
	gorm.Model
	OwnerId    uint
	PublicId   string
	BusinessId uint
	Points     uint

	Transactions []Transaction
	OwnedItems   []OwnedItem
}

type Transaction struct {
	gorm.Model
	PublicId      string
	VirtualCardId uint
	Code          string
	State         TransactionStateEnum
	AddedPoints   uint

	TransactionDetails []TransactionDetail
}

type TransactionDetail struct {
	gorm.Model
	TransactionId uint
	ItemId        uint
	Action        ActionTypeEnum
}
