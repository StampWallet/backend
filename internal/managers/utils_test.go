package managers

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/lithammer/shortuuid/v4"
	"gorm.io/gorm"

	. "github.com/StampWallet/backend/internal/database"
	//. "github.com/StampWallet/backend/internal/database/mocks"
	. "github.com/StampWallet/backend/internal/services/mocks"
)

func getTestUser(db GormDB) *User {
    user := User {
	PublicId: "test",
	FirstName: "first",
	LastName: "last",
	Email: "test@example.com",
	PasswordHash: "test",
	EmailVerified: true,
    }
    tx := db.Create(&user)
    if err := tx.GetError(); err != nil {
	panic(fmt.Errorf("failed to create Item %w", err))
    }
    return &user
}

func getTestBusiness(db GormDB, user *User) *Business {
    business := Business{
	PublicId: shortuuid.New(),
	OwnerId: user.ID,
	Name: "test business",
	Description: "Description",
	Address: "test address",
	GPSCoordinates: "+27.5916+086.5640+8850CRSWGS_84/",
	NIP: "1234567890",
	KRS: "1234567890",
	REGON: "1234567890",
	OwnerName: "test owner",
    }
    tx := db.Create(&business)
    if err := tx.GetError(); err != nil {
	panic(fmt.Errorf("failed to create Business %w", err))
    }
    return &business
}

func MockStorage(user *User, storage *MockFileStorageService) FileMetadata {
    file := FileMetadata {
	Model: gorm.Model {
	    ID: uint(rand.Uint32()),
	},
	PublicId: shortuuid.New(),
	OwnerId: user.ID,
    }
    storage.
	EXPECT().
	CreateStub(&user).
	Return(file)
    return file
}

func getTestItemDefinition(db GormDB, business *Business, metadata FileMetadata) *ItemDefinition {
    definition := ItemDefinition{
	PublicId: shortuuid.New(),
	BusinessId: business.ID,
	Name: "test item definition name",
	Price: 10,
	Description: "test item definition description",
	ImageId: metadata.PublicId,
	StartDate: time.Now(),
	EndDate: time.Now().Add(time.Hour*24),
	MaxAmount: 10,
	Available: true,
    }
    tx := db.Create(&definition)
    if err := tx.GetError(); err != nil {
	panic(fmt.Errorf("failed to create ItemDefinition %w", err))
    }
    return &definition
}

func getTestVirtualCard(db GormDB, user *User, business *Business) *VirtualCard {
    virtualCard := VirtualCard {
	PublicId: shortuuid.New(),
	OwnerId: user.ID,
	BusinessId: business.ID,
	Points: 40,
    }
    tx := db.Create(&virtualCard)
    if err := tx.GetError(); err != nil {
	panic(fmt.Errorf("failed to create VirtualCard %w", err))
    }
    return &virtualCard
}

func getTestOwnedItem(db GormDB, itemDefinition *ItemDefinition, card *VirtualCard) *OwnedItem {
    ownedItem := OwnedItem {
	PublicId: shortuuid.New(),
	DefinitionId: itemDefinition.ID,
	VirtualCardId: card.ID,
	Used: sql.NullTime{ Valid: false },
	Status: OwnedItemStatusOwned,
    }
    tx := db.Create(&ownedItem)
    if err := tx.GetError(); err != nil {
	panic(fmt.Errorf("failed to create OwnedItem %w", err))
    }
    return &ownedItem
}

func getTestLocalCard(db GormDB, user *User) *LocalCard {
    localCard := LocalCard {
	PublicId: shortuuid.New(),
	OwnerId: user.ID,
	Type: CardTypes[0].PublicId,
	Code: "012345678901",
	Name: "test card",
    }
    tx := db.Create(&localCard)
    if err := tx.GetError(); err != nil {
	panic(fmt.Errorf("failed to create LocalCard %w", err))
    }
    return &localCard
}
