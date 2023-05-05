package testutils

import (
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/lithammer/shortuuid/v4"
	"gorm.io/gorm"

	. "github.com/StampWallet/backend/internal/database"
	//. "github.com/StampWallet/backend/internal/database/mocks"
	//. "github.com/StampWallet/backend/internal/services/mocks"
)

func Save[T any](db GormDB, item *T) {
	if db == nil {
		return
	}
	tx := db.Create(item)
	if err := tx.GetError(); err != nil {
		panic(fmt.Errorf("failed to create item of type %T: %w", *new(T), err))
	}
}

func GetTestUser(db GormDB) *User {
	user := User{
		PublicId:      shortuuid.New(),
		FirstName:     shortuuid.New(),
		LastName:      shortuuid.New(),
		Email:         shortuuid.New() + "@example.com",
		PasswordHash:  shortuuid.New(),
		EmailVerified: true,
	}
	Save(db, &user)
	return &user
}

func GetDefaultUser() *User {
	return GetTestUser(nil)
}

func GetTestBusiness(db GormDB, user *User) *Business {
	business := Business{
		PublicId:       shortuuid.New(),
		OwnerId:        user.ID,
		Name:           "test business",
		Description:    "Description",
		Address:        "test address",
		GPSCoordinates: FromCoords(27.5916, 086.5640),
		NIP:            strconv.Itoa(rand.Intn(math.MaxInt)), // TODO: this needs to generate valid codes for tests
		KRS:            strconv.Itoa(rand.Intn(math.MaxInt)),
		REGON:          strconv.Itoa(rand.Intn(math.MaxInt)),
		OwnerName:      "test owner",
		BannerImageId:  GetTestFileMetadata(db, user).PublicId,
		IconImageId:    GetTestFileMetadata(db, user).PublicId,
	}
	Save(db, &business)
	return &business
}

func GetDefaultBusiness(user *User) *Business {
	return GetTestBusiness(nil, user)
}

func GetTestFileMetadata(db GormDB, user *User) *FileMetadata {
	file := FileMetadata{
		Model: gorm.Model{
			ID: uint(rand.Uint32()),
		},
		PublicId: shortuuid.New(),
		OwnerId:  user.ID,
	}
	Save(db, &file)
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
		StartDate:   sql.NullTime{Time: time.Now(), Valid: true},
		EndDate:     sql.NullTime{Time: time.Now().Add(time.Hour * 24), Valid: true},
		MaxAmount:   10,
		Available:   true,
	}
	Save(db, &definition)
	return &definition
}

func GetTestVirtualCard(db GormDB, user *User, business *Business) *VirtualCard {
	virtualCard := VirtualCard{
		PublicId:   shortuuid.New(),
		OwnerId:    user.ID,
		BusinessId: business.ID,
		Points:     40,
	}
	Save(db, &virtualCard)
	return &virtualCard
}

func GetTestOwnedItem(db GormDB, itemDefinition *ItemDefinition, card *VirtualCard) *OwnedItem {
	ownedItem := OwnedItem{
		PublicId:      shortuuid.New(),
		DefinitionId:  itemDefinition.ID,
		VirtualCardId: card.ID,
		Used:          sql.NullTime{Valid: false},
		Status:        OwnedItemStatusOwned,
	}
	Save(db, &ownedItem)
	return &ownedItem
}

func GetDefaultOwnedItem(itemDefinition *ItemDefinition, card *VirtualCard) *OwnedItem {
	return GetTestOwnedItem(nil, itemDefinition, card)
}

func GetTestLocalCard(db GormDB, user *User) *LocalCard {
	localCard := LocalCard{
		PublicId: shortuuid.New(),
		OwnerId:  user.ID,
		Type:     "test type",
		Code:     strconv.Itoa(rand.Intn(math.MaxInt)),
		Name:     "test card",
	}
	Save(db, &localCard)
	return &localCard
}

func GetTestTransaction(db GormDB, virtualCard *VirtualCard, items []OwnedItem) (*Transaction, []TransactionDetail) {
	transaction := Transaction{
		PublicId:      shortuuid.New(),
		VirtualCardId: virtualCard.ID,
		Code:          strconv.Itoa(rand.Intn(math.MaxInt)),
		State:         TransactionStateStarted,
		AddedPoints:   0,
	}
	Save(db, &transaction)
	var details []TransactionDetail
	for _, item := range items {
		transactionDetail := TransactionDetail{
			TransactionId: transaction.ID,
			ItemId:        item.ID,
			Action:        NoActionType,
		}
		details = append(details, transactionDetail)
		Save(db, &transactionDetail)
	}
	return &transaction, details
}

func GetDefaultItem(business *Business) *ItemDefinition {
	return &ItemDefinition{
		PublicId:    shortuuid.New(),
		BusinessId:  business.ID,
		Name:        "test item definition name",
		Price:       10,
		Description: "test item definition description",
		ImageId:     "does not matter",
		StartDate:   sql.NullTime{Time: time.Now(), Valid: true},
		EndDate:     sql.NullTime{Time: time.Now().Add(time.Hour * 24), Valid: true},
		MaxAmount:   10,
		Available:   true,
	}
}
