package managers

import (
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

func GetItemDefinitionManager(ctrl *gomock.Controller) *ItemDefinitionManagerImpl {
	return &ItemDefinitionManagerImpl{
		&BaseServices{
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
		CreateStub(&user).
		Return(*imageFile, nil)

	details := &ItemDetails{
		Name:        "test item",
		Price:       Ptr(uint(50)),
		Description: "item description",
		StartDate:   Ptr(time.Now()),
		EndDate:     Ptr(time.Now().Add(time.Hour * 24)),
		MaxAmount:   Ptr(uint(10)),
		Available:   Ptr(true),
	}
	definition, err := manager.AddItem(business, details)
	require.Truef(t, MatchEntities(details, definition), "entities do not match")
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
	manager.fileStorageService.(*MockFileStorageService).
		EXPECT().
		CreateStub(&user).
		Return(*imageFile, nil)

	newDetails := ItemDetails{
		Name:        "new item details",
		Price:       Ptr(uint(10)),
		Description: "new item description",
		StartDate:   Ptr(time.Now().Add(time.Hour * 24)),
		EndDate:     Ptr(time.Now().Add(time.Hour * 48)),
		MaxAmount:   Ptr(uint(20)),
		Available:   Ptr(false),
	}
	newDefinition, err := manager.ChangeItemDetails(definition, &newDetails)

	require.Truef(t, MatchEntities(newDetails, newDefinition), "entities do not match")
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

	manager.fileStorageService.(*MockFileStorageService).
		EXPECT().
		CreateStub(&user).
		Return(*imageFile, nil)
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
	require.Equalf(t, OwnedItemStatusWithdrawn, dbOwnedItem.Status, "new owned item status is not withdrawn")
	require.Equalf(t, virtualCard.Points+definition.Price, dbVirtualCard.Points, "db virtual card did not regain points")
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
	manager.fileStorageService.(*MockFileStorageService).
		EXPECT().
		CreateStub(&user).
		Return(*iconImage, nil)

	returnedDefinitions, err := manager.GetForBusiness(business)
	require.Nilf(t, err, "GetForBusiness returned an error")
	require.Equalf(t, 1, len(returnedDefinitions), "GetForBusiness returned more or less than one result")
	require.Equalf(t, returnedDefinitions[0].PublicId, definition.PublicId, "GetForBusiness returned definition name does not match whats expected")
}
