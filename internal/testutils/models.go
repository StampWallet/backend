package testutils

import (
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/lithammer/shortuuid/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	. "github.com/StampWallet/backend/internal/database"
	//. "github.com/StampWallet/backend/internal/database/mocks"
	//. "github.com/StampWallet/backend/internal/services/mocks"
)

func Save(db GormDB, item any) {
	if db == nil {
		return
	}
	tx := db.Session(&gorm.Session{FullSaveAssociations: true}).Save(item)
	//fmt.Printf("after save: %v\n", item)
	if err := tx.GetError(); err != nil {
		panic(fmt.Errorf("failed to create item of type %T: %w", item, err))
	}
}

func GetTestUser(db GormDB) *User {
	user := User{
		PublicId: shortuuid.New(),
		//FirstName:     shortuuid.New(),
		//LastName:      shortuuid.New(),
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
	a := GetTestFileMetadata(db, user).PublicId
	b := GetTestFileMetadata(db, user).PublicId
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
		BannerImageId:  a,
		IconImageId:    b,
		User:           user,
	}
	Save(db, &business)
	user.Business = &business
	business.User = user
	return &business
}

func GetDefaultBusiness(user *User) *Business {
	return GetTestBusiness(nil, user)
}

func GetTestFileMetadata(db GormDB, user *User) *FileMetadata {
	file := FileMetadata{
		PublicId: shortuuid.New(),
		OwnerId:  user.ID,
		User:     user,
	}
	Save(db, &user)
	Save(db, &file)
	file.User = user
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
		Business:    business,
	}
	Save(db, &definition)
	definition.Business = business
	return &definition
}

func GetTestItemDefinitionWithPrice(db GormDB, business *Business, metadata FileMetadata, price uint) *ItemDefinition {
	definition := ItemDefinition{
		PublicId:    shortuuid.New(),
		BusinessId:  business.ID,
		Name:        "test item definition name",
		Price:       price,
		Description: "test item definition description",
		ImageId:     metadata.PublicId,
		StartDate:   sql.NullTime{Time: time.Now(), Valid: true},
		EndDate:     sql.NullTime{Time: time.Now().Add(time.Hour * 24), Valid: true},
		MaxAmount:   10,
		Available:   true,
		Business:    business,
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
		User:       user,
		Business:   business,
	}
	Save(db, &virtualCard)
	return &virtualCard
}

func GetTestVirtualCardWithPoints(db GormDB, user *User, business *Business, points uint) *VirtualCard {
	virtualCard := VirtualCard{
		PublicId:   shortuuid.New(),
		OwnerId:    user.ID,
		BusinessId: business.ID,
		Points:     points,
		User:       user,
		Business:   business,
	}
	Save(db, &virtualCard)
	return &virtualCard
}

func GetTestOwnedItem(db GormDB, itemDefinition *ItemDefinition, card *VirtualCard) *OwnedItem {
	ownedItem := OwnedItem{
		PublicId:       shortuuid.New(),
		DefinitionId:   itemDefinition.ID,
		VirtualCardId:  card.ID,
		Used:           sql.NullTime{Valid: false},
		Status:         OwnedItemStatusOwned,
		ItemDefinition: itemDefinition,
		VirtualCard:    card,
	}
	Save(db, &ownedItem)
	return &ownedItem
}

func GetTestOwnedItemUsed(db GormDB, itemDefinition *ItemDefinition, card *VirtualCard) *OwnedItem {
	ownedItem := OwnedItem{
		PublicId:       shortuuid.New(),
		DefinitionId:   itemDefinition.ID,
		VirtualCardId:  card.ID,
		Used:           sql.NullTime{Valid: true, Time: time.Now()},
		Status:         OwnedItemStatusUsed,
		VirtualCard:    card,
		ItemDefinition: itemDefinition,
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
			OwnedItem:     &item,
		}
		details = append(details, transactionDetail)
		Save(db, &transactionDetail)
	}
	transaction.TransactionDetails = details
	return &transaction, details
}

func GetDefaultItem(business *Business) *ItemDefinition {
	itemDefinition := &ItemDefinition{
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
		Business:    business,
	}
	return itemDefinition
}

func GetTestToken(db GormDB, user *User) (*Token, string) {
	secret := shortuuid.New()
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), 10)
	if err != nil {
		panic(fmt.Sprintf("failed to create token hash %s", err))
	}
	token := Token{
		TokenId:      shortuuid.New(),
		TokenHash:    string(hash),
		Expires:      time.Now().Add(24 * time.Hour),
		TokenPurpose: TokenPurposeEmail,
		Used:         false,
		Recalled:     false,
		OwnerId:      user.ID,
		User:         user,
	}
	Save(db, &token)
	token.User = user
	return &token, secret
}

func GetTestSessionToken(db GormDB, user *User, expires time.Time) (*Token, string) {
	secret := shortuuid.New()
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), 10)
	if err != nil {
		panic(fmt.Sprintf("failed to create token hash %s", err))
	}
	token := Token{
		TokenId:      shortuuid.New(),
		TokenHash:    string(hash),
		Expires:      expires,
		TokenPurpose: TokenPurposeSession,
		Used:         false,
		Recalled:     false,
		User:         user,
	}
	Save(db, &token)
	token.User = user
	return &token, secret
}

func GetTestMenuImage(db GormDB, business *Business) *MenuImage {
	fileId := GetTestFileMetadata(db, business.User).PublicId
	menuImage := MenuImage{
		FileId:     fileId,
		BusinessId: business.ID,
		Business:   business,
	}
	Save(db, &menuImage)
	return &menuImage
}
