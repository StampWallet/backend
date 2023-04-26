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

func getItemDefinitionManager(ctrl *gomock.Controller) *ItemDefinitionManagerImpl {
    return &ItemDefinitionManagerImpl {
        &BaseServices {
            Logger: log.Default(),
	    Database: GetDatabase(),
        },
	NewMockFileStorageService(ctrl),
    }
}

func TestItemDefinitionAddItem(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager := getItemDefinitionManager(ctrl)
    user := getTestUser(manager.baseServices.Database)
    business := getTestBusiness(manager.baseServices.Database, user)
    imageFile := MockStorage(user, manager.fileStorageService.(*MockFileStorageService))
    details := &ItemDetails {
	Name: "test item",
	Price: Ptr(uint(50)),
	Description: "item description",
	StartDate: Ptr(time.Now()),
	EndDate: Ptr(time.Now().Add(time.Hour*24)),
	MaxAmount: Ptr(uint(10)),
	Available: Ptr(true),
    }
    definition, err := manager.AddItem(business, details)
    require.Truef(t, MatchEntities(details, definition), "entities do not match")
    require.Equalf(t, err, nil, "additem retuned an error")
    require.Equalf(t, definition.ImageId, imageFile.PublicId, "additem retuned an error")
    var dbDetails ItemDefinition
    tx := manager.baseServices.Database.Find(&dbDetails, ItemDefinition { Model: gorm.Model { ID: definition.ID } })
    require.Equalf(t, tx.GetError(), nil, "database find returned an error")
    require.Equalf(t, dbDetails.Name, definition.Name, "database find returned invalid data")
}

func TestItemDefinitionChangeItemDetails(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager := getItemDefinitionManager(ctrl)
    user := getTestUser(manager.baseServices.Database)
    business := getTestBusiness(manager.baseServices.Database, user)
    definition := getTestItemDefinition(manager.baseServices.Database, business, 
	MockStorage(user, manager.fileStorageService.(*MockFileStorageService)))

    newDetails := ItemDetails{
	Name: "new item details",
	Price: Ptr(uint(10)),
	Description: "new item description",
	StartDate: Ptr(time.Now().Add(time.Hour*24)),
	EndDate: Ptr(time.Now().Add(time.Hour*48)),
	MaxAmount: Ptr(uint(20)),
	Available: Ptr(false),
    }
    newDefinition, err := manager.ChangeItemDetails(definition, &newDetails)

    require.Truef(t, MatchEntities(newDetails, newDefinition), "entities do not match")
    require.Equalf(t, err, nil, "ChangeItemDetails retuned an error")

    var dbDetails ItemDefinition
    tx := manager.baseServices.Database.Find(&dbDetails, ItemDefinition { Model: gorm.Model { ID: definition.ID } })
    require.Equalf(t, tx.GetError(), nil, "database find returned an error")
    require.Equalf(t, dbDetails.Name, newDetails.Name, "database find returned invalid data")
}

func TestItemDefinitionWithdrawItem(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager := getItemDefinitionManager(ctrl)
    user := getTestUser(manager.baseServices.Database)
    business := getTestBusiness(manager.baseServices.Database, user)
    definition := getTestItemDefinition(manager.baseServices.Database, business, 
	MockStorage(user, manager.fileStorageService.(*MockFileStorageService)))
    virtualCard := getTestVirtualCard(manager.baseServices.Database, user, business)
    ownedItem := getTestOwnedItem(manager.baseServices.Database, definition, virtualCard)

    newDefinition, err := manager.WithdrawItem(definition)
    require.Equalf(t, err, nil, "WithdrawItem returned an error")

    var dbItemDefinition ItemDefinition
    var dbOwnedItem OwnedItem
    var dbVirtualCard VirtualCard

    tx := manager.baseServices.Database.Find(&dbItemDefinition, ItemDefinition { Model: gorm.Model { ID: definition.ID } })
    require.Equalf(t, tx.GetError(), nil, "db ItemDefinition find returned an error")
    tx = manager.baseServices.Database.Find(&dbOwnedItem, OwnedItem { Model: gorm.Model { ID: ownedItem.ID } })
    require.Equalf(t, tx.GetError(), nil, "db OwnedItem find returned an error")
    tx = manager.baseServices.Database.Find(&dbVirtualCard, VirtualCard { Model: gorm.Model { ID: virtualCard.ID } })
    require.Equalf(t, tx.GetError(), nil, "db VirtualCard find returned an error")

    require.Equalf(t, dbItemDefinition.Withdrawn, true, "db item definition is not withdrawn")
    require.Equalf(t, newDefinition.Withdrawn, true, "new item definition is not withdrawn")
    require.Equalf(t, dbOwnedItem.Status, OwnedItemStatusWithdrawn, "new owned item status is not withdrawn")
    require.Equalf(t, dbVirtualCard.Points, virtualCard.Points + definition.Price, "db virtual card did not regain points")
}

func TestItemDefinitionGetForBusiness(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager := getItemDefinitionManager(ctrl)
    user := getTestUser(manager.baseServices.Database)
    business := getTestBusiness(manager.baseServices.Database, user)
    definition := getTestItemDefinition(manager.baseServices.Database, business, 
	MockStorage(user, manager.fileStorageService.(*MockFileStorageService)))

    returnedDefinitions, err := manager.GetForBusiness(business)
    require.Equalf(t, err, nil, "GetForBusiness returned an error")
    require.Equalf(t, len(returnedDefinitions), 1, "GetForBusiness returned more or less than one result")
    require.Equalf(t, returnedDefinitions[0].PublicId, definition.PublicId, "GetForBusiness returned definition name does not match whats expected")
}
