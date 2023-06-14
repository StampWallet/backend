package managers

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	. "github.com/StampWallet/backend/internal/database"
	//. "github.com/StampWallet/backend/internal/database/mocks"
	. "github.com/StampWallet/backend/internal/services"
	. "github.com/StampWallet/backend/internal/services/mocks"
	. "github.com/StampWallet/backend/internal/testutils"
)

type ItemDetailsMatcher struct {
	Name        string
	Price       *uint
	Description string
	StartDate   sql.NullTime
	EndDate     sql.NullTime
	MaxAmount   *uint
	Available   *bool
}

func GetItemDefinitionManager(ctrl *gomock.Controller) *ItemDefinitionManagerImpl {
	return &ItemDefinitionManagerImpl{
		BaseServices{
			Logger:   log.Default(),
			Database: GetTestDatabase(),
		},
		NewMockFileStorageService(ctrl),
	}
}

func TestItemDefinitionAddItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetItemDefinitionManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	business := GetTestBusiness(manager.baseServices.Database, user)
	imageFile := GetTestFileMetadata(manager.baseServices.Database, user)
	manager.fileStorageService.(*MockFileStorageService).
		EXPECT().
		CreateStub(user).
		Return(imageFile, nil)

	details := &ItemDetails{
		Name:        "test item",
		Price:       Ptr(uint(50)),
		Description: "item description",
		StartDate:   Ptr(time.Now()),
		EndDate:     Ptr(time.Now().Add(time.Hour * 24)),
		MaxAmount:   Ptr(uint(10)),
		Available:   Ptr(true),
	}
	detailsMatcher := &ItemDetailsMatcher{
		Name:        details.Name,
		Price:       details.Price,
		Description: details.Description,
		StartDate:   sql.NullTime{*details.StartDate, true},
		EndDate:     sql.NullTime{*details.EndDate, true},
		MaxAmount:   details.MaxAmount,
		Available:   details.Available,
	}
	definition, err := manager.AddItem(user, business, details)
	require.Truef(t, MatchEntities(detailsMatcher, definition), "entities do not match")
	require.Nilf(t, err, "additem returned an error")
	require.Equalf(t, imageFile.PublicId, definition.ImageId, "additem returned an error")
	var dbDetails ItemDefinition
	tx := manager.baseServices.Database.Find(&dbDetails, ItemDefinition{Model: gorm.Model{ID: definition.ID}})
	require.Nilf(t, tx.GetError(), "database find returned an error")
	require.Equalf(t, definition.Name, dbDetails.Name, "database find returned invalid data")
}

func TestItemDefinitionChangeItemDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetItemDefinitionManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	business := GetTestBusiness(manager.baseServices.Database, user)
	imageFile := GetTestFileMetadata(manager.baseServices.Database, user)
	definition := GetTestItemDefinition(manager.baseServices.Database, business, *imageFile)

	newDetails := ItemDetails{
		Name:        "new item details",
		Price:       Ptr(uint(10)),
		Description: "new item description",
		StartDate:   Ptr(time.Now().Add(time.Hour * 24)),
		EndDate:     Ptr(time.Now().Add(time.Hour * 48)),
		MaxAmount:   Ptr(uint(20)),
		Available:   Ptr(false),
	}
	detailsMatcher := &ItemDetailsMatcher{
		Name:        newDetails.Name,
		Price:       newDetails.Price,
		Description: newDetails.Description,
		StartDate:   sql.NullTime{*newDetails.StartDate, true},
		EndDate:     sql.NullTime{*newDetails.EndDate, true},
		MaxAmount:   newDetails.MaxAmount,
		Available:   newDetails.Available,
	}
	newDefinition, err := manager.ChangeItemDetails(definition, &newDetails)

	require.Truef(t, MatchEntities(detailsMatcher, newDefinition), "entities do not match")
	require.Nilf(t, err, "ChangeItemDetails returned an error")

	var dbDetails ItemDefinition
	tx := manager.baseServices.Database.Find(&dbDetails, ItemDefinition{Model: gorm.Model{ID: definition.ID}})
	require.Nilf(t, tx.GetError(), "database find returned an error")
	require.Equalf(t, newDetails.Name, dbDetails.Name, "database find returned invalid data")
}

func TestItemDefinitionWithdrawItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetItemDefinitionManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	business := GetTestBusiness(manager.baseServices.Database, user)
	imageFile := GetTestFileMetadata(manager.baseServices.Database, user)
	definition := GetTestItemDefinition(manager.baseServices.Database, business, *imageFile)

	// manager.fileStorageService.(*MockFileStorageService).
	// 	EXPECT().
	// 	CreateStub(&user).
	// 	Return(imageFile, nil) // __jm__ why?
	virtualCard := GetTestVirtualCard(manager.baseServices.Database, user, business)
	ownedItem := GetTestOwnedItem(manager.baseServices.Database, definition, virtualCard)

	newDefinition, err := manager.WithdrawItem(definition)
	require.Nilf(t, err, "WithdrawItem returned an error")

	var dbItemDefinition ItemDefinition
	var dbOwnedItem OwnedItem
	var dbVirtualCard VirtualCard

	tx := manager.baseServices.Database.Find(&dbItemDefinition, ItemDefinition{Model: gorm.Model{ID: definition.ID}})
	require.Nilf(t, tx.GetError(), "db ItemDefinition find returned an error")
	tx = manager.baseServices.Database.Find(&dbOwnedItem, OwnedItem{Model: gorm.Model{ID: ownedItem.ID}})
	require.Nilf(t, tx.GetError(), "db OwnedItem find returned an error")
	tx = manager.baseServices.Database.Find(&dbVirtualCard, VirtualCard{Model: gorm.Model{ID: virtualCard.ID}})
	require.Nilf(t, tx.GetError(), "db VirtualCard find returned an error")

	require.Equalf(t, true, dbItemDefinition.Withdrawn, "db item definition is not withdrawn")
	require.Equalf(t, true, newDefinition.Withdrawn, "new item definition is not withdrawn")
	require.Equalf(t, OwnedItemStatusEnum(OwnedItemStatusWithdrawn), dbOwnedItem.Status, "new owned item status is not withdrawn")
	require.Equalf(t, virtualCard.Points+definition.Price, dbVirtualCard.Points, "db virtual card did not regain points")
}

func TestItemDefinitionWithdrawItemMultiple(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetItemDefinitionManager(ctrl)
	db := manager.baseServices.Database

	//create two business users
	userB1 := GetTestUser(manager.baseServices.Database)
	userB2 := GetTestUser(manager.baseServices.Database)

	// create two businesses
	business1 := GetTestBusiness(manager.baseServices.Database, userB1)
	business2 := GetTestBusiness(manager.baseServices.Database, userB2)

	// create two item definitions per business
	definition1B1 := GetTestItemDefinitionWithPrice(manager.baseServices.Database, business1,
		*GetTestFileMetadata(manager.baseServices.Database, userB1), 10)
	definition2B1 := GetTestItemDefinitionWithPrice(manager.baseServices.Database, business1,
		*GetTestFileMetadata(manager.baseServices.Database, userB1), 20)
	definition1B2 := GetTestItemDefinitionWithPrice(manager.baseServices.Database, business1,
		*GetTestFileMetadata(manager.baseServices.Database, userB2), 30)
	definition2B2 := GetTestItemDefinitionWithPrice(manager.baseServices.Database, business1,
		*GetTestFileMetadata(manager.baseServices.Database, userB2), 40)

	// create two users
	user1 := GetTestUser(manager.baseServices.Database)
	user2 := GetTestUser(manager.baseServices.Database)

	// create two virtual cards for each user and business
	virtualCardU1B1 := GetTestVirtualCardWithPoints(manager.baseServices.Database, user1, business1, 0)
	virtualCardU1B2 := GetTestVirtualCardWithPoints(manager.baseServices.Database, user1, business2, 10)
	virtualCardU2B1 := GetTestVirtualCardWithPoints(manager.baseServices.Database, user2, business1, 20)
	virtualCardU2B2 := GetTestVirtualCardWithPoints(manager.baseServices.Database, user2, business2, 30)

	// create three owned items per item definition per virtual card, one used
	var ownedItemsToBeWithdrawn []*OwnedItem = []*OwnedItem{
		GetTestOwnedItem(manager.baseServices.Database, definition1B1, virtualCardU1B1), // 10 points
		GetTestOwnedItem(manager.baseServices.Database, definition1B1, virtualCardU1B1), // 10 points

		GetTestOwnedItem(manager.baseServices.Database, definition1B1, virtualCardU2B1), // 10 points
		GetTestOwnedItem(manager.baseServices.Database, definition1B1, virtualCardU2B1), // 10 points
	}

	var ownedItemsToRemainUnmodified []*OwnedItem = []*OwnedItem{
		GetTestOwnedItemUsed(manager.baseServices.Database, definition1B1, virtualCardU1B1),

		GetTestOwnedItemUsed(manager.baseServices.Database, definition1B1, virtualCardU2B1),

		GetTestOwnedItem(manager.baseServices.Database, definition2B1, virtualCardU1B1),
		GetTestOwnedItem(manager.baseServices.Database, definition2B1, virtualCardU1B1),

		GetTestOwnedItem(manager.baseServices.Database, definition1B2, virtualCardU1B2),
		GetTestOwnedItem(manager.baseServices.Database, definition1B2, virtualCardU1B2),
		GetTestOwnedItem(manager.baseServices.Database, definition2B2, virtualCardU1B2),
		GetTestOwnedItem(manager.baseServices.Database, definition2B2, virtualCardU1B2),

		GetTestOwnedItem(manager.baseServices.Database, definition2B1, virtualCardU2B1),
		GetTestOwnedItem(manager.baseServices.Database, definition2B1, virtualCardU2B1),

		GetTestOwnedItem(manager.baseServices.Database, definition1B2, virtualCardU2B2),
		GetTestOwnedItem(manager.baseServices.Database, definition1B2, virtualCardU2B2),
		GetTestOwnedItem(manager.baseServices.Database, definition2B2, virtualCardU2B2),
		GetTestOwnedItem(manager.baseServices.Database, definition2B2, virtualCardU2B2),
	}

	// withdraw one item definition

	itemDefinition, err := manager.WithdrawItem(definition1B1)
	require.Nilf(t, err, "manager.WithdrawItem returned an error")
	require.Equalf(t, true, itemDefinition.Withdrawn, "manager.WithdrawItem returned an item that was not  withdrawn")

	for _, item := range ownedItemsToBeWithdrawn {
		var dbItem OwnedItem
		tx := db.Find(&dbItem, OwnedItem{Model: gorm.Model{ID: item.Model.ID}})
		require.Nilf(t, tx.GetError(), "db.Find returned an error")
		require.Equalf(t, OwnedItemStatusEnum(OwnedItemStatusWithdrawn), dbItem.Status,
			"item was not withdrawn")
	}

	for _, item := range ownedItemsToRemainUnmodified {
		var dbItem OwnedItem
		tx := db.Find(&dbItem, OwnedItem{Model: gorm.Model{ID: item.Model.ID}})
		require.Nilf(t, tx.GetError(), "db.Find returned an error")
		require.Equalf(t, item.Status, dbItem.Status, "item was modified")
	}

	// check if result is expected
	//should be modified (+20 points each)
	for _, card := range []*VirtualCard{virtualCardU1B1, virtualCardU2B1} {
		var dbCard VirtualCard
		tx := db.Find(&dbCard, VirtualCard{Model: gorm.Model{ID: card.Model.ID}})
		require.Nilf(t, tx.GetError(), "db.Find returned an error")
		require.Equalf(t, card.Points+20, dbCard.Points, "card has a different amount of points than expected")
	}

	//shouldnt be modified
	for _, card := range []*VirtualCard{virtualCardU1B2, virtualCardU2B2} {
		var dbCard VirtualCard
		tx := db.Find(&dbCard, VirtualCard{Model: gorm.Model{ID: card.Model.ID}})
		require.Nilf(t, tx.GetError(), "db.Find returned an error")
		require.Equalf(t, card.Points, dbCard.Points, "card was modified")
	}
}

func TestItemDefinitionGetForBusiness(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetItemDefinitionManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	business := GetTestBusiness(manager.baseServices.Database, user)
	iconImage := GetTestFileMetadata(manager.baseServices.Database, user)
	definition := GetTestItemDefinition(manager.baseServices.Database, business,
		*iconImage)

	returnedDefinitions, err := manager.GetForBusiness(business)
	require.Nilf(t, err, "GetForBusiness returned an error")
	require.Equalf(t, 1, len(returnedDefinitions), "GetForBusiness returned more or less than one result")
	require.Equalf(t, returnedDefinitions[0].PublicId, definition.PublicId, "GetForBusiness returned definition name does not match whats expected")
}
