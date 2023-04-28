package testutils

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/lithammer/shortuuid/v4"
	"gorm.io/gorm"

	. "github.com/StampWallet/backend/internal/database"
	//. "github.com/StampWallet/backend/internal/database/mocks"
	//. "github.com/StampWallet/backend/internal/services/mocks"
)

func GetDefaultUser() *User {
	return &User{
		PublicId:      "test",
		FirstName:     "first",
		LastName:      "last",
		Email:         "test@example.com",
		PasswordHash:  "test",
		EmailVerified: true,
	}
}

func GetTestUser(db GormDB) *User {
	userPtr := GetDefaultUser()
	tx := db.Create(userPtr)
	if err := tx.GetError(); err != nil {
		panic(fmt.Errorf("failed to create Item %w", err))
	}
	return userPtr
}

func GetTestBusiness(db GormDB, user *User) *Business {
	business := Business{
		PublicId:       shortuuid.New(),
		OwnerId:        user.ID,
		Name:           "test business",
		Description:    "Description",
		Address:        "test address",
		GPSCoordinates: "+27.5916+086.5640+8850CRSWGS_84/",
		NIP:            "1234567890",
		KRS:            "1234567890",
		REGON:          "1234567890",
		OwnerName:      "test owner",
	}
	tx := db.Create(&business)
	if err := tx.GetError(); err != nil {
		panic(fmt.Errorf("failed to create Business %w", err))
	}
	return &business
}

func GetTestFileMetadata(db GormDB, user *User) *FileMetadata {
	file := FileMetadata{
		Model: gorm.Model{
			ID: uint(rand.Uint32()),
		},
		PublicId: shortuuid.New(),
		OwnerId:  user.ID,
	}
	tx := db.Create(&file)
	if err := tx.GetError(); err != nil {
		panic(fmt.Errorf("failed to create Business %w", err))
	}
	return &file
}

func GetTestItemDefinition(db GormDB, business *Business, metadata FileMetadata) *ItemDefinition {
	definition := ItemDefinition{
		PublicId:    shortuuid.New(),
		BusinessId:  business.ID,
		Name:        "test item definition name",
		Price:       10,
		Description: "test item definition description",
		ImageId:     metadata.PublicId,
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour * 24),
		MaxAmount:   10,
		Available:   true,
	}
	tx := db.Create(&definition)
	if err := tx.GetError(); err != nil {
		panic(fmt.Errorf("failed to create ItemDefinition %w", err))
	}
	return &definition
}

func GetTestVirtualCard(db GormDB, user *User, business *Business) *VirtualCard {
	virtualCard := VirtualCard{
		PublicId:   shortuuid.New(),
		OwnerId:    user.ID,
		BusinessId: business.ID,
		Points:     40,
	}
	tx := db.Create(&virtualCard)
	if err := tx.GetError(); err != nil {
		panic(fmt.Errorf("failed to create VirtualCard %w", err))
	}
	return &virtualCard
}

func GetDefaultOwnedItem(itemDefinition *ItemDefinition, card *VirtualCard) *OwnedItem {
	return &OwnedItem{
		PublicId:      shortuuid.New(),
		DefinitionId:  itemDefinition.ID,
		VirtualCardId: card.ID,
		Used:          sql.NullTime{Valid: false},
		Status:        OwnedItemStatusOwned,
	}
}

func Save(db GormDB, item any) {
	tx := db.Create(item)
	if err := tx.GetError(); err != nil {
		panic(fmt.Errorf("failed to create OwnedItem %w", err))
	}
}

func GetTestOwnedItem(db GormDB, itemDefinition *ItemDefinition, card *VirtualCard) *OwnedItem {
	ownedItem := GetDefaultOwnedItem(itemDefinition, card)
	Save(db, ownedItem)
	return ownedItem
}

func GetTestLocalCard(db GormDB, user *User) *LocalCard {
	localCard := LocalCard{
		PublicId: shortuuid.New(),
		OwnerId:  user.ID,
		Type:     "test type",
		Code:     "012345678901",
		Name:     "test card",
	}
	tx := db.Create(&localCard)
	if err := tx.GetError(); err != nil {
		panic(fmt.Errorf("failed to create LocalCard %w", err))
	}
	return &localCard
}

func GetTestTransaction(db GormDB, virtualCard *VirtualCard, items []OwnedItem) (*Transaction, []TransactionDetail) {
	transaction := Transaction{
		PublicId:      shortuuid.New(),
		VirtualCardId: virtualCard.ID,
		Code:          strconv.Itoa(rand.Intn(100)),
		State:         TransactionStateStarted,
		AddedPoints:   0,
	}
	tx := db.Create(&transaction)
	if err := tx.GetError(); err != nil {
		panic(fmt.Errorf("failed to create Transaction %w", err))
	}
	var details []TransactionDetail
	for _, item := range items {
		transactionDetail := TransactionDetail{
			TransactionId: transaction.ID,
			ItemId:        item.ID,
			Action:        NoActionType,
		}
		tx := db.Create(&transactionDetail)
		details = append(details, transactionDetail)
		if err := tx.GetError(); err != nil {
			panic(fmt.Errorf("failed to create TransactionDetail %w", err))
		}
	}
	return &transaction, details
}

func GetDefaultItem(business *Business) ItemDefinition {
	return ItemDefinition{
		PublicId:    shortuuid.New(),
		BusinessId:  business.ID,
		Name:        "test item definition name",
		Price:       10,
		Description: "test item definition description",
		ImageId:     "does not matter",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour * 24),
		MaxAmount:   10,
		Available:   true,
	}
}

func SaveItem(db GormDB, definition *ItemDefinition) {
	tx := db.Create(&definition)
	if err := tx.GetError(); err != nil {
		panic(fmt.Errorf("failed to create ItemDefinition %w", err))
	}
}
