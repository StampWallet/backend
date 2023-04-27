package managers

import (
	"log"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services"
	. "github.com/StampWallet/backend/internal/testutils"
)

func GetTestVirtualCardManager(ctrl *gomock.Controller) *VirtualCardManagerImpl {
	return &VirtualCardManagerImpl{
		&BaseServices{
			Logger:   log.Default(),
			Database: GetDatabase(),
		},
	}
}

type virtualCardManagerTest struct {
	ctrl           *gomock.Controller
	manager        *VirtualCardManagerImpl
	user           *User
	businessUser   *User
	business       *Business
	itemDefinition *ItemDefinition
	db             GormDB
}

type businessTest struct {
	businessUser   *User
	business       *Business
	itemDefinition *ItemDefinition
}

func setupBusiness(db GormDB) businessTest {
	businessUser := GetTestUser(db)
	business := GetTestBusiness(db, businessUser)
	itemDefinition := GetTestItemDefinition(db, business, *GetTestFileMetadata(db, businessUser))
	return businessTest{
		businessUser:   businessUser,
		business:       business,
		itemDefinition: itemDefinition,
	}
}

func setupVirualCardManagerTest(t *testing.T) virtualCardManagerTest {
	ctrl := gomock.NewController(t)
	manager := GetTestVirtualCardManager(gomock.NewController(t))
	db := manager.baseServices.Database
	user := GetTestUser(db)
	businessTest := setupBusiness(db)
	return virtualCardManagerTest{
		ctrl:           ctrl,
		manager:        manager,
		user:           user,
		businessUser:   businessTest.businessUser,
		business:       businessTest.business,
		itemDefinition: businessTest.itemDefinition,
		db:             db,
	}
}

func TestVirtualCardManagerCreate(t *testing.T) {
	s := setupVirualCardManagerTest(t)
	virtualCard, err := s.manager.Create(*s.user, *s.business)
	require.Nilf(t, err, "VirtualCardManager.Create should retrun a nil error")
	require.NotNilf(t, virtualCard, "VirtualCardManager.Create shuold not return a nil virtual card")
	if virtualCard == nil {
		return
	}
	require.Equalf(t, s.user.ID, virtualCard.OwnerId, "VirtualCardManager.Create shuold returned a card that belongs to the passed user")
	require.Equalf(t, s.business.ID, virtualCard.BusinessId, "VirtualCardManager.Create should return a card that belongs to the passed business")
	require.Equalf(t, 0, virtualCard.Points, "VirtualCardManager.Create should returned a card with 0 points")

	newVirtualCard, newErr := s.manager.Create(*s.user, *s.business)
	require.Equalf(t, VirtualCardAlreadyExists, newErr, "VirtualCardManager.Create should returned VirtualCardAlreadyExists if the user attempts to create the same card twice")
	require.Nilf(t, newVirtualCard, "VirtualCardManager.Create should return a nil pointer if the user attempts to create the same card twice")
}

func TestVirtualCardManagerRemove(t *testing.T) {
	s := setupVirualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	ownedItem := GetTestOwnedItem(s.db, s.itemDefinition, virtualCard)

	err := s.manager.Remove(virtualCard)
	require.Nilf(t, err, "VirtualCardManager.Remove should not retrun a nil error")

	var dbVirtualCard VirtualCard
	var dbOwnedItem OwnedItem
	tx := s.db.First(&dbVirtualCard, VirtualCard{Model: gorm.Model{ID: virtualCard.ID}})
	require.Equalf(t, gorm.ErrRecordNotFound, tx.GetError(), "Datbase find for VirtualCard should return RecordNotFound")
	tx = s.db.First(&dbOwnedItem, OwnedItem{Model: gorm.Model{ID: ownedItem.ID}})
	require.Equalf(t, gorm.ErrRecordNotFound, tx.GetError(), "Datbase find for OwnedItem should return RecordNotFound")
}

func TestVirtualCardManagerRemoveNotExisting(t *testing.T) {
	s := setupVirualCardManagerTest(t)

	err := s.manager.Remove(&VirtualCard{Model: gorm.Model{ID: 123123}})
	require.Equalf(t, NoSuchVirtualCard, err, "VirtualCardManager.Remove should retrun NoSuchVirtualCard error")
}

func TestVirtualCardManagerGetForUser(t *testing.T) {
	s := setupVirualCardManagerTest(t)
	cards, err := s.manager.GetForUser(s.user)
	require.Nilf(t, err, "VirtualCardManager.GetForUser should return a nil error")
	require.Equalf(t, 0, len(cards), "VirtualCardManager.GetForUser should return zero cards")

	anotherBusiness := setupBusiness(s.db)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	GetTestOwnedItem(s.db, s.itemDefinition, virtualCard)
	virtualCard2 := GetTestVirtualCard(s.db, s.user, anotherBusiness.business)
	GetTestOwnedItem(s.db, s.itemDefinition, virtualCard)

	//random virtual card that should not be present in resuls
	GetTestVirtualCard(s.db, GetTestUser(s.db), s.business)

	cards, err = s.manager.GetForUser(s.user)
	require.Nilf(t, err, "VirtualCardManager.GetForUser should return a nil error")
	require.Equalf(t, 2, len(cards), "VirtualCardManager.GetForUser should return two cards")
	nCards := 0
	for _, c := range cards {
		if c.PublicId == virtualCard.PublicId || c.PublicId == virtualCard2.PublicId {
			nCards += 1
		} else {
			require.Failf(t, "VirtualCardManager.GetForUser returned an unexpected card %s", c.PublicId)
		}
	}
	require.Equalf(t, 2, nCards, "VirtualCardManager.GetForUser should return two cards")
}

func TestVirtualCardManagerGetOwnedItems(t *testing.T) {
	s := setupVirualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	ownedItem := GetTestOwnedItem(s.db, s.itemDefinition, virtualCard)

	items, err := s.manager.GetOwnedItems(virtualCard)
	require.Nilf(t, err, "VirtualCardManager.GetOwnedItems should return a nil error")
	require.Equalf(t, 1, len(items), "VirtualCardManager.GetOwnedItems should return one item")
	require.Equalf(t, ownedItem.PublicId, items[0].PublicId, "VirtualCardManager.GetOwnedItems should return the expected item")
}

func TestVirtualCardManagerFilterOwnedItems(t *testing.T) {
	s := setupVirualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	ownedItem := GetTestOwnedItem(s.db, s.itemDefinition, virtualCard)
	GetTestOwnedItem(s.db, s.itemDefinition, virtualCard)
	ownedItem3 := GetTestOwnedItem(s.db, s.itemDefinition, virtualCard)

	//random virtualcard/user/item, should not appear in the results
	GetTestOwnedItem(s.db, s.itemDefinition, GetTestVirtualCard(s.db, GetTestUser(s.db), s.business))

	items, err := s.manager.FilterOwnedItems(virtualCard, []string{ownedItem.PublicId, ownedItem3.PublicId})
	require.Nilf(t, err, "VirtualCardManager.FilterOwnedItems should return a nil error")
	require.Equalf(t, 2, len(items), "VirtualCardManager.FilterOwnedItems should return two items")
	nItems := 0
	for _, i := range items {
		if i.PublicId == ownedItem.PublicId || i.PublicId == ownedItem3.PublicId {
			nItems += 1
		} else {
			require.Failf(t, "VirtualCardManager.FilterOwnedItems returned an unexpected item %s", i.PublicId)
		}
	}
	require.Equalf(t, 2, nItems, "VirtualCardManager.FilterOwnedItems should return two items")
}

func TestVirtualCardManagerBuyItem(t *testing.T) {
	s := setupVirualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	itemDefinition := GetDefaultItem(s.business)
	SaveItem(s.db, &itemDefinition)

	ownedItem, err := s.manager.BuyItem(virtualCard, &itemDefinition)
	require.Nilf(t, err, "VirtualCardManager.BuyItem should return a nil error")
	require.NotNilf(t, ownedItem, "VirtualCardManager.BuyItem should not return a nil ownedItem")
	if virtualCard == nil {
		return
	}
	require.Equalf(t, itemDefinition.ID, ownedItem.DefinitionId, "VirtualCardManager.BuyItem should return item from expected item definition")
	require.Equalf(t, OwnedItemStatusOwned, ownedItem.Status, "VirtualCardManager.BuyItem should return item from expected item definition")
	require.Falsef(t, ownedItem.Used.Valid, "VirtualCardManager.BuyItem should return item with nil used time")
	require.Equalf(t, virtualCard.ID, ownedItem.VirtualCardId, "VirtualCardManager.BuyItem should return item owned by expected virtual card")

	var dbVirtualCard VirtualCard
	tx := s.db.First(&dbVirtualCard, VirtualCard{Model: gorm.Model{ID: virtualCard.ID}})
	require.Nilf(t, tx.GetError(), "Database find should not return an error")
	require.Truef(t, MatchEntities(virtualCard, dbVirtualCard), "Virtual card should be the same as in the db")
}

func TestVirtualCardManagerBuyItemAboveMaxAmount(t *testing.T) {
	s := setupVirualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	itemDefinition := GetDefaultItem(s.business)
	itemDefinition.MaxAmount = 1
	SaveItem(s.db, &itemDefinition)

	ownedItem, err := s.manager.BuyItem(virtualCard, &itemDefinition)
	require.NotNilf(t, ownedItem, "VirtualCardManager.BuyItem should not return a nil item")
	require.Nilf(t, err, "VirtualCardManager.BuyItem should return a nil error")

	ownedItemAboveMaxAmount, err := s.manager.BuyItem(virtualCard, &itemDefinition)
	require.Nilf(t, ownedItemAboveMaxAmount, "VirtualCardManager.BuyItem should return a nil item")
	require.Equalf(t, err, AboveMaxAmount, "VirtualCardManager.BuyItem should return a AboveMaxAmount error")
}

func TestVirtualCardManagerBuyItemWithNotEnoughPoints(t *testing.T) {
	s := setupVirualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	itemDefinition := GetDefaultItem(s.business)
	itemDefinition.Price = virtualCard.Points + 1
	SaveItem(s.db, &itemDefinition)

	ownedItem, err := s.manager.BuyItem(virtualCard, &itemDefinition)
	require.Nilf(t, ownedItem, "VirtualCardManager.BuyItem should return a nil item")
	require.Equalf(t, NotEnoughPoints, err, "VirtualCardManager.BuyItem should a NotEnoughPoints error")
}

func TestVirtualCardManagerBuyItemBeforeStartDate(t *testing.T) {
	s := setupVirualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	itemDefinition := GetDefaultItem(s.business)
	itemDefinition.StartDate = time.Now().Add(time.Hour * 24 * 24)
	itemDefinition.EndDate = time.Now().Add(time.Hour * 24 * 25)
	SaveItem(s.db, &itemDefinition)

	ownedItem, err := s.manager.BuyItem(virtualCard, &itemDefinition)
	require.Nilf(t, ownedItem, "VirtualCardManager.BuyItem should return a nil item")
	require.Equalf(t, BeforeStartDate, err, "VirtualCardManager.BuyItem should return a BeforeStartDate error")
}

func TestVirtualCardManagerBuyItemAfterEndDate(t *testing.T) {
	s := setupVirualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	itemDefinition := GetDefaultItem(s.business)
	itemDefinition.StartDate = time.Now().Add(-time.Hour * 24 * 24)
	itemDefinition.EndDate = time.Now().Add(-time.Hour)
	SaveItem(s.db, &itemDefinition)

	ownedItem, err := s.manager.BuyItem(virtualCard, &itemDefinition)
	require.Nilf(t, ownedItem, "VirtualCardManager.BuyItem should return a nil item")
	require.Equalf(t, AfterEndDate, err, "VirtualCardManager.BuyItem should return an AfterEndDate error")
}

func TestVirtualCardManagerBuyItemUnavailable(t *testing.T) {
	s := setupVirualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	itemDefinition := GetDefaultItem(s.business)
	itemDefinition.Available = false
	SaveItem(s.db, &itemDefinition)

	ownedItem, err := s.manager.BuyItem(virtualCard, &itemDefinition)
	require.Nilf(t, ownedItem, "VirtualCardManager.BuyItem should return a nil item")
	require.Equalf(t, UnavailableItem, err, "VirtualCardManager.BuyItem should return an UnavailableItem error")
}

func TestVirtualCardManagerBuyItemWithdrawn(t *testing.T) {
	s := setupVirualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	itemDefinition := GetDefaultItem(s.business)
	itemDefinition.Withdrawn = false
	SaveItem(s.db, &itemDefinition)

	ownedItem, err := s.manager.BuyItem(virtualCard, &itemDefinition)
	require.Nilf(t, ownedItem, "VirtualCardManager.BuyItem should return a nil item")
	require.Equalf(t, WithdrawnItem, err, "VirtualCardManager.BuyItem should return a WithdrawnItem error")
}

func TestVirtualCardManagerReturnItem(t *testing.T) {
	s := setupVirualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	ownedItem := GetTestOwnedItem(s.db, s.itemDefinition, virtualCard)

	err := s.manager.ReturnItem(ownedItem)
	require.Nilf(t, err, "VirtualCardManager.Return should return a nil error")

	var dbVirtualCard VirtualCard
	tx := s.db.First(&dbVirtualCard, VirtualCard{Model: gorm.Model{ID: virtualCard.ID}})
	require.Nilf(t, tx.GetError(), "Datbase find for OwnedItem should not return an error")
	require.Equalf(t, virtualCard.Points+s.itemDefinition.Price, dbVirtualCard.Points, "Virtual card points amount should be updated")

	var dbOwnedItem OwnedItem
	tx = s.db.First(&dbOwnedItem, OwnedItem{Model: gorm.Model{ID: virtualCard.ID}})
	require.Nilf(t, tx.GetError(), "Datbase find for OwnedItem should not return an error")
	require.Equalf(t, OwnedItemStatusReturned, dbOwnedItem.Status, "Owned item in the database should have returned status")

	err = s.manager.ReturnItem(ownedItem)
	require.Equalf(t, ItemReturned, err, "VirtualCardManager.Return should return an ItemReturned error on second return try")

	tx = s.db.First(&dbVirtualCard, VirtualCard{Model: gorm.Model{ID: virtualCard.ID}})
	require.Nilf(t, tx.GetError(), "Datbase find for OwnedItem should not return an error")
	require.Equalf(t, virtualCard.Points+s.itemDefinition.Price, dbVirtualCard.Points, "Virtual card points amount should stay the same on second return try")
}
