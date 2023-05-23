package database

import (
	"context"
	"database/sql"
	//"database/sql/driver"
	"fmt"
	"time"

	"github.com/twpayne/go-geom"
	//"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	PublicId string `gorm:"uniqueIndex;not null"`
	OwnerId  uint
	Type     string
	Code     string
	Name     string

	User *User `gorm:"foreignkey:OwnerId"`
}

func (entity *LocalCard) GetUserId(_ GormDB) (uint, error) {
	return entity.OwnerId, nil
}

// Token

type Token struct {
	gorm.Model
	OwnerId      uint
	TokenId      string           `gorm:"uniqueIndex;not null"`
	TokenHash    string           `gorm:"not null"`
	Expires      time.Time        `gorm:"not null"`
	TokenPurpose TokenPurposeEnum `gorm:"not null"`
	Used         bool             `gorm:"default:false;not null"`
	Recalled     bool             `gorm:"default:false;not null"`

	User *User `gorm:"foreignkey:OwnerId"`
}

func (entity *Token) GetUserId(_ GormDB) (uint, error) {
	return entity.OwnerId, nil
}

// FileMetadata

type FileMetadata struct {
	gorm.Model
	PublicId    string `gorm:"uniqueIndex;not null"`
	OwnerId     uint   `gorm:"not null"`
	ContentType sql.NullString
	Uploaded    sql.NullTime

	User *User `gorm:"foreignkey:OwnerId"`
}

func (entity *FileMetadata) GetUserId(_ GormDB) (uint, error) {
	return entity.OwnerId, nil
}

// User

type User struct {
	gorm.Model
	PublicId string `gorm:"uniqueIndex;not null"`
	//FirstName     string `gorm:"not null"`
	//LastName      string `gorm:"not null"`
	Email         string `gorm:"uniqueIndex;not null"`
	PasswordHash  string `gorm:"not null"`
	EmailVerified bool   `gorm:"default:false;not null"`

	Tokens        []Token        `gorm:"foreignkey:OwnerId"`
	FilesMetadata []FileMetadata `gorm:"foreignkey:OwnerId"`
	LocalCards    []LocalCard    `gorm:"foreignkey:OwnerId"`
	VirtualCards  []VirtualCard  `gorm:"foreignkey:OwnerId"`
	Business      *Business      `gorm:"foreignkey:OwnerId"`
}

// Business

type GPSCoordinates geom.Point

func (g *GPSCoordinates) Scan(input interface{}) error {
	gt, err := ewkbhex.Decode(input.(string))
	if err != nil {
		return err
	}
	gp := gt.(*geom.Point)
	gc := GPSCoordinates(*gp)
	g = &gc
	return nil
}

//// TODO probably does not work
//func (g GPSCoordinates) Value() (driver.Value, error) {
//	b := geom.Point(g)
//	//return ewkbhex.Encode(&b, ewkb.NDR)
//	//fmt.Printf("%f %f\n", b.X(), b.Y())
//	if b.Empty() {
//		return "SRID=3857;POINT(0 0)", nil
//	} else {
//		return fmt.Sprintf("SRID=3857;POINT(%f %f)", b.X(), b.Y()), nil
//	}
//}

func (g GPSCoordinates) GormDataType() string {
	return "geometry"
}

func (g GPSCoordinates) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	b := geom.Point(g)
	var vars []interface{} = []interface{}{"SRID=4326;POINT(0 0)"}
	if !b.Empty() {
		vars = []interface{}{fmt.Sprintf("SRID=4326;POINT(%f %f)", b.X(), b.Y())}
	}
	return clause.Expr{
		SQL:  "ST_PointFromText(?)",
		Vars: vars,
	}
}

func FromCoords(longitude float64, latitude float64) GPSCoordinates {
	return GPSCoordinates(*geom.NewPointFlat(geom.XY, geom.Coord{longitude, latitude}))
}

// NOTE index is created in automigrate. couldn't figure out how to create a gin fulltext on many columns from tags alone.
// INDEX NAME is business_fulltext_idx
type Business struct {
	gorm.Model
	PublicId       string         `gorm:"uniqueIndex;not null"`
	OwnerId        uint           `gorm:"not null"`
	Name           string         `gorm:"not null"`
	Description    string         `gorm:"not null"`
	Address        string         `gorm:"not null"`
	GPSCoordinates GPSCoordinates `gorm:"type:geography(POINT,4326);index:,type:gist"`
	NIP            string         `gorm:"unique;not null"`
	KRS            string         `gorm:"unique;not null"`
	REGON          string         `gorm:"unique;not null"`
	OwnerName      string         `gorm:"not null"`
	BannerImageId  string         `gorm:"unique;not null"`
	IconImageId    string         `gorm:"unique;not null"`

	ItemDefinitions []ItemDefinition `gorm:"foreignkey:BusinessId"`
	MenuItems       []MenuItem       `gorm:"foreignkey:BusinessId"`
	VirtualCards    []VirtualCard    `gorm:"foreignkey:BusinessId"`

	User *User `gorm:"foreignkey:OwnerId"`
}

func (entity *Business) GetUserId() (uint, error) {
	return entity.OwnerId, nil
}

// ItemDefinition

type ItemDefinition struct {
	gorm.Model
	PublicId    string `gorm:"uniqueIndex;not null"`
	BusinessId  uint   `gorm:"not null"`
	Name        string `gorm:"not null"`
	Price       uint   `gorm:"not null"`
	Description string `gorm:"not null"`
	ImageId     string `gorm:"unique"`
	StartDate   sql.NullTime
	EndDate     sql.NullTime
	MaxAmount   uint
	Available   bool `gorm:"not null"`
	Withdrawn   bool `gorm:"not null"`

	OwnedItems []OwnedItem `gorm:"foreignkey:DefinitionId"`

	Business *Business `gorm:"foreignkey:BusinessId"`
}

func (entity *ItemDefinition) GetBusinessId(_ GormDB) (uint, error) {
	return entity.BusinessId, nil
}

// MenuItem

type MenuItem struct {
	gorm.Model
	BusinessId uint   `gorm:"uniqueIndex:menu_item_idx;not null"`
	FileId     string `gorm:"uniqueIndex:menu_item_idx;not null"`

	Business *Business `gorm:"foreignkey:BusinessId"`
}

func (entity *MenuItem) GetBusinessId() (uint, error) {
	return entity.BusinessId, nil
}

// OwnedItem

type OwnedItem struct {
	gorm.Model
	PublicId      string `gorm:"uniqueIndex;not null"`
	DefinitionId  uint   `gorm:"not null"`
	VirtualCardId uint   `gorm:"not null"`
	Used          sql.NullTime
	Status        OwnedItemStatusEnum `gorm:"default:OWNED;not null"`

	ItemDefinition *ItemDefinition `gorm:"foreignkey:DefinitionId"`
	VirtualCard    *VirtualCard    `gorm:"foreignkey:VirtualCardId"`
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
	PublicId   string `gorm:"uniqueIndex;not null"`
	OwnerId    uint   `gorm:"not null"`
	BusinessId uint   `gorm:"not null"`
	Points     uint   `gorm:"not null"`

	OwnedItems   []OwnedItem   `gorm:"foreignkey:VirtualCardId"`
	Transactions []Transaction `gorm:"foreignkey:VirtualCardId"`

	Business *Business `gorm:"foreignkey:BusinessId"`
	User     *User     `gorm:"foreignkey:OwnerId"`
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
	PublicId      string               `gorm:"uniqueIndex;not null"`
	VirtualCardId uint                 `gorm:"index:code,unique,priority:1;not null"`
	Code          string               `gorm:"index:code,unique,priority:2;not null"`
	State         TransactionStateEnum `gorm:"default:STARTED;not null"`
	AddedPoints   uint

	TransactionDetails []TransactionDetail `gorm:"foreignkey:TransactionId"`

	VirtualCard *VirtualCard `gorm:"foreignkey:VirtualCardId"`
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
	TransactionId uint           `gorm:"index:transaction_detail,unique,priority:1;not null"`
	ItemId        uint           `gorm:"index:transaction_detail,unique,priority:2;not null"`
	Action        ActionTypeEnum `gorm:"default:NO_ACTION;not null"`

	Transaction *Transaction `gorm:"foreignkey:TransactionId"`
}
