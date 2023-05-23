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
	. "github.com/StampWallet/backend/internal/services"
	. "github.com/StampWallet/backend/internal/testutils"
)

// Builds VirtualCardManagerImpl with test database
func GetTestVirtualCardManager(ctrl *gomock.Controller) *VirtualCardManagerImpl {
	return &VirtualCardManagerImpl{
		&BaseServices{
			Logger:   log.Default(),
			Database: GetDatabase(),
		},
	}
}

// Return type of setupVirtualCardManagerTes
type virtualCardManagerTest struct {
	ctrl           *gomock.Controller
	manager        *VirtualCardManagerImpl
	user           *User
	businessUser   *User
	business       *Business
	itemDefinition *ItemDefinition
	db             GormDB
}

// Return type of setupBusiness
type businessTest struct {
	businessUser   *User
	business       *Business
	itemDefinition *ItemDefinition
}

// Creates and example business, with a new owner and example ItemDefinition
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

// Creates a new user, business, virtualCard and VirtualCardManagerImpl
func setupVirtualCardManagerTest(t *testing.T) virtualCardManagerTest {
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

// Tests VirtualCardManagerImpl.Create on happy path and when virtualCard for business and user already exists
func TestVirtualCardManagerCreate(t *testing.T) {
	s := setupVirtualCardManagerTest(t)
	virtualCard, err := s.manager.Create(*s.user, *s.business)
	require.Nilf(t, err, "VirtualCardManager.Create should return a nil error")
	require.NotNilf(t, virtualCard, "VirtualCardManager.Create should not return a nil virtual card")
	if virtualCard == nil {
		return
	}
	require.Equalf(t, s.user.ID, virtualCard.OwnerId, "VirtualCardManager.Create should returned a card that belongs to the passed user")
	require.Equalf(t, s.business.ID, virtualCard.BusinessId, "VirtualCardManager.Create should return a card that belongs to the passed business")
	require.Equalf(t, uint(0), virtualCard.Points, "VirtualCardManager.Create should returned a card with 0 points")

	newVirtualCard, newErr := s.manager.Create(*s.user, *s.business)
	require.Equalf(t, ErrVirtualCardAlreadyExists, newErr, "VirtualCardManager.Create should returned VirtualCardAlreadyExists if the user attempts to create the same card twice")
	require.Nilf(t, newVirtualCard, "VirtualCardManager.Create should return a nil pointer if the user attempts to create the same card twice")
}

// Tests VirtualCardManagerImpl.Remove on happy path
func TestVirtualCardManagerRemove(t *testing.T) {
	s := setupVirtualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	ownedItem := GetTestOwnedItem(s.db, s.itemDefinition, virtualCard)

	user2 := GetTestUser(s.db)
	virtualCard2 := GetTestVirtualCard(s.db, user2, s.business)
	ownedItem2 := GetTestOwnedItem(s.db, s.itemDefinition, virtualCard2)

	err := s.manager.Remove(virtualCard)
	require.Nilf(t, err, "VirtualCardManager.Remove should retrun a nil error")

	var dbVirtualCard VirtualCard
	var dbOwnedItem OwnedItem
	tx := s.db.First(&dbVirtualCard, VirtualCard{Model: gorm.Model{ID: virtualCard.ID}})
	require.Equalf(t, gorm.ErrRecordNotFound, tx.GetError(), "Datbase find for VirtualCard should return RecordNotFound")
	tx = s.db.First(&dbOwnedItem, OwnedItem{Model: gorm.Model{ID: ownedItem.ID}})
	require.Equalf(t, gorm.ErrRecordNotFound, tx.GetError(), "Datbase find for OwnedItem should return RecordNotFound")

	tx = s.db.First(&dbVirtualCard, VirtualCard{Model: gorm.Model{ID: virtualCard2.ID}})
	require.Equalf(t, virtualCard2.PublicId, dbVirtualCard.PublicId, "Datbase find for VirtualCard2 should return VirtualCard2")
	tx = s.db.First(&dbOwnedItem, OwnedItem{Model: gorm.Model{ID: ownedItem2.ID}})
	require.Equalf(t, ownedItem2.PublicId, dbOwnedItem.PublicId, "Datbase find for OwnedItem2 should return OwnedItem2")
}

// Tests VirtualCardManagerImpl.Remove when virtualCard does not exist
func TestVirtualCardManagerRemoveNotExisting(t *testing.T) {
	s := setupVirtualCardManagerTest(t)

	err := s.manager.Remove(&VirtualCard{Model: gorm.Model{ID: 123123}})
	require.Equalf(t, ErrNoSuchVirtualCard, err, "VirtualCardManager.Remove should retrun NoSuchVirtualCard error")
}

// Tests VirtualCardManagerImpl.GetOwnedItems on happy path
func TestVirtualCardManagerGetOwnedItems(t *testing.T) {
	s := setupVirtualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	ownedItem := GetTestOwnedItem(s.db, s.itemDefinition, virtualCard)

	items, err := s.manager.GetOwnedItems(virtualCard)
	require.Nilf(t, err, "VirtualCardManager.GetOwnedItems should return a nil error")
	require.Equalf(t, 1, len(items), "VirtualCardManager.GetOwnedItems should return one item")
	require.Equalf(t, ownedItem.PublicId, items[0].PublicId, "VirtualCardManager.GetOwnedItems should return the expected item")
}

// Tests VirtualCardManagerImpl.FilterOwnedItems on happy path
func TestVirtualCardManagerFilterOwnedItems(t *testing.T) {
	s := setupVirtualCardManagerTest(t)
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

// Tests VirtualCardManagerImpl.BuyItem on happy path
func TestVirtualCardManagerBuyItem(t *testing.T) {
	s := setupVirtualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	itemDefinition := GetDefaultItem(s.business)
	Save(s.db, itemDefinition)
	oldPoints := virtualCard.Points

	ownedItem, err := s.manager.BuyItem(virtualCard, itemDefinition)
	require.Nilf(t, err, "VirtualCardManager.BuyItem should return a nil error")
	require.NotNilf(t, ownedItem, "VirtualCardManager.BuyItem should not return a nil ownedItem")
	require.Equalf(t, itemDefinition.ID, ownedItem.DefinitionId, "VirtualCardManager.BuyItem should return item from expected item definition")
	require.Equalf(t, OwnedItemStatusOwned, ownedItem.Status, "VirtualCardManager.BuyItem should return item from expected item definition")
	require.Falsef(t, ownedItem.Used.Valid, "VirtualCardManager.BuyItem should return item with nil used time")
	require.Equalf(t, virtualCard.ID, ownedItem.VirtualCardId, "VirtualCardManager.BuyItem should return item owned by expected virtual card")

	var dbVirtualCard VirtualCard
	tx := s.db.First(&dbVirtualCard, VirtualCard{Model: gorm.Model{ID: virtualCard.ID}})
	require.Nilf(t, tx.GetError(), "Database find should not return an error")
	require.Equalf(t, oldPoints-itemDefinition.Price, dbVirtualCard.Points,
		"Virtual card from the db should have its points substracted")
	require.Equalf(t, oldPoints-itemDefinition.Price, virtualCard.Points,
		"Virtual card from args should have its points substracted")
}

// Tests VirtualCardManagerImpl.BuyItem when virtualCard has ItemDefinition.MaxAmount items
func TestVirtualCardManagerBuyItemAboveMaxAmount(t *testing.T) {
	s := setupVirtualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	itemDefinition := GetDefaultItem(s.business)
	itemDefinition.MaxAmount = 1
	Save(s.db, itemDefinition)

	ownedItem, err := s.manager.BuyItem(virtualCard, itemDefinition)
	require.NotNilf(t, ownedItem, "VirtualCardManager.BuyItem should not return a nil item")
	require.Nilf(t, err, "VirtualCardManager.BuyItem should return a nil error")

	ownedItemAboveMaxAmount, err := s.manager.BuyItem(virtualCard, itemDefinition)
	require.Nilf(t, ownedItemAboveMaxAmount, "VirtualCardManager.BuyItem should return a nil item")
	require.Equalf(t, err, ErrAboveMaxAmount, "VirtualCardManager.BuyItem should return a AboveMaxAmount error")
}

// Tests VirtualCardManagerImpl.BuyItem when virtualCard does not have enough points
func TestVirtualCardManagerBuyItemWithNotEnoughPoints(t *testing.T) {
	s := setupVirtualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	itemDefinition := GetDefaultItem(s.business)
	itemDefinition.Price = virtualCard.Points + 1
	Save(s.db, itemDefinition)

	ownedItem, err := s.manager.BuyItem(virtualCard, itemDefinition)
	require.Nilf(t, ownedItem, "VirtualCardManager.BuyItem should return a nil item")
	require.Equalf(t, ErrNotEnoughPoints, err, "VirtualCardManager.BuyItem should a NotEnoughPoints error")
}

// Tests VirtualCardManagerImpl.BuyItem when ItemDefinition is not valid yet
func TestVirtualCardManagerBuyItemBeforeStartDate(t *testing.T) {
	s := setupVirtualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	itemDefinition := GetDefaultItem(s.business)
	itemDefinition.StartDate = sql.NullTime{Time: time.Now().Add(time.Hour * 24 * 24), Valid: true}
	itemDefinition.EndDate = sql.NullTime{Time: time.Now().Add(time.Hour * 24 * 25), Valid: true}
	Save(s.db, itemDefinition)

	ownedItem, err := s.manager.BuyItem(virtualCard, itemDefinition)
	require.Nilf(t, ownedItem, "VirtualCardManager.BuyItem should return a nil item")
	require.Equalf(t, ErrBeforeStartDate, err, "VirtualCardManager.BuyItem should return a BeforeStartDate error")
}

// Tests VirtualCardManagerImpl.BuyItem when ItemDefinition is already invalid
func TestVirtualCardManagerBuyItemAfterEndDate(t *testing.T) {
	s := setupVirtualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	itemDefinition := GetDefaultItem(s.business)
	itemDefinition.StartDate = sql.NullTime{Time: time.Now().Add(-time.Hour * 24 * 24), Valid: true}
	itemDefinition.EndDate = sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true}
	Save(s.db, itemDefinition)

	ownedItem, err := s.manager.BuyItem(virtualCard, itemDefinition)
	require.Nilf(t, ownedItem, "VirtualCardManager.BuyItem should return a nil item")
	require.Equalf(t, ErrAfterEndDate, err, "VirtualCardManager.BuyItem should return an AfterEndDate error")
}

// Tests VirtualCardManagerImpl.BuyItem when ItemDefinition is unavailable
func TestVirtualCardManagerBuyItemUnavailable(t *testing.T) {
	s := setupVirtualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	itemDefinition := GetDefaultItem(s.business)
	itemDefinition.Available = false
	Save(s.db, itemDefinition)

	ownedItem, err := s.manager.BuyItem(virtualCard, itemDefinition)
	require.Nilf(t, ownedItem, "VirtualCardManager.BuyItem should return a nil item")
	require.Equalf(t, ErrUnavailableItem, err, "VirtualCardManager.BuyItem should return an UnavailableItem error")
}

// Tests VirtualCardManagerImpl.BuyItem when ItemDefinition is withdrawn
func TestVirtualCardManagerBuyItemWithdrawn(t *testing.T) {
	s := setupVirtualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	itemDefinition := GetDefaultItem(s.business)
	itemDefinition.Withdrawn = true
	Save(s.db, itemDefinition)

	ownedItem, err := s.manager.BuyItem(virtualCard, itemDefinition)
	require.Nilf(t, ownedItem, "VirtualCardManager.BuyItem should return a nil item")
	require.Equalf(t, ErrWithdrawnItem, err, "VirtualCardManager.BuyItem should return a WithdrawnItem error")
}

// Tests [VirtualCardManagerImpl.ReturnItem] on happy path and when item was already returned
func TestVirtualCardManagerReturnItem(t *testing.T) {
	s := setupVirtualCardManagerTest(t)
	virtualCard := GetTestVirtualCard(s.db, s.user, s.business)
	ownedItem := GetTestOwnedItem(s.db, s.itemDefinition, virtualCard)

	err := s.manager.ReturnItem(ownedItem)
	require.Nilf(t, err, "VirtualCardManager.Return should return a nil error")

	var dbVirtualCard VirtualCard
	tx := s.db.First(&dbVirtualCard, VirtualCard{Model: gorm.Model{ID: virtualCard.ID}})
	require.Nilf(t, tx.GetError(), "Datbase find for OwnedItem should not return an error")
	require.Equalf(t, virtualCard.Points+s.itemDefinition.Price, dbVirtualCard.Points, "Virtual card points amount should be updated")

	var dbOwnedItem OwnedItem
	tx = s.db.First(&dbOwnedItem, OwnedItem{Model: gorm.Model{ID: ownedItem.ID}})
	require.Nilf(t, tx.GetError(), "Datbase find for OwnedItem should not return an error")
	require.Equalf(t, OwnedItemStatusEnum(OwnedItemStatusReturned), dbOwnedItem.Status, "Owned item in the database should have returned status")

	err = s.manager.ReturnItem(ownedItem)
	require.Equalf(t, ErrItemCantBeReturned, err, "VirtualCardManager.Return should return an ItemCantBeReturned error on second return try")

	tx = s.db.First(&dbVirtualCard, VirtualCard{Model: gorm.Model{ID: virtualCard.ID}})
	require.Nilf(t, tx.GetError(), "Datbase find for OwnedItem should not return an error")
	require.Equalf(t, virtualCard.Points+s.itemDefinition.Price, dbVirtualCard.Points, "Virtual card points amount should stay the same on second return try")
}
