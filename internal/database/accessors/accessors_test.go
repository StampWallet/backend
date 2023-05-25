package database

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	//"github.com/stretchr/testify/require"
	//"gorm.io/gorm"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/testutils"
)

func GetBusinessAuthorizedAccessor(ctrl *gomock.Controller) *BusinessAuthorizedAccessorImpl {
	return &BusinessAuthorizedAccessorImpl{
		database: GetTestDatabase(),
	}
}

func GetUserAuthorizedAccessor(ctrl *gomock.Controller) *UserAuthorizedAccessorImpl {
	return &UserAuthorizedAccessorImpl{
		database: GetTestDatabase(),
	}
}

func GetTransactionAuthorizedAccessor(ctrl *gomock.Controller) *AuthorizedTransactionAccessorImpl {
	return &AuthorizedTransactionAccessorImpl{
		database: GetTestDatabase(),
	}
}

// BusinessAuthorizedAccessor

func setupBusinessAccessorTest(t *testing.T) (*gomock.Controller, *BusinessAuthorizedAccessorImpl, *User, *Business, *ItemDefinition) {
	ctrl := gomock.NewController(t)
	accessor := GetBusinessAuthorizedAccessor(ctrl)
	user := GetTestUser(accessor.database)
	business := GetTestBusiness(accessor.database, user)
	itemDefinition := GetTestItemDefinition(accessor.database, business,
		*GetTestFileMetadata(accessor.database, user))
	return ctrl, accessor, user, business, itemDefinition
}

func TestBusinessAuthorizedAccessorValidAccess(t *testing.T) {
	_, accessor, _, business, itemDefinition := setupBusinessAccessorTest(t)
	result, err := accessor.Get(business, &ItemDefinition{PublicId: itemDefinition.PublicId})
	require.NotNilf(t, result, "accessor returned nil")
	require.Nil(t, err, "accessor returned non nil error")
	obtainedItemDefinition := result.(*ItemDefinition)
	require.Equal(t, itemDefinition.PublicId, obtainedItemDefinition.PublicId, "accessor returned unexpected entity")
}

func TestBusinessAuthorizedAccessorInvalidAccess(t *testing.T) {
	_, accessor, user, _, itemDefinition := setupBusinessAccessorTest(t)
	anotherBusiness := GetTestBusiness(accessor.database, user)

	result, err := accessor.Get(anotherBusiness, &ItemDefinition{PublicId: itemDefinition.PublicId})
	require.Nil(t, result, "accessor did not return nil")
	require.Equalf(t, ErrNoAccess, err, "accessor returned error other than NoAccess")
}

func TestBusinessAuthorizedAccessorNoEntity(t *testing.T) {
	_, accessor, _, business, _ := setupBusinessAccessorTest(t)

	result, err := accessor.Get(business, &ItemDefinition{PublicId: "test"})
	require.Nil(t, result, "accessor did not return nil")
	require.Equal(t, ErrNotFound, err, "accessor returned error other than NoAccess")
}

// UserAuthorizedAccessor

func setupUserAccessorTest(t *testing.T) (*gomock.Controller, *UserAuthorizedAccessorImpl, *User, *LocalCard) {
	ctrl := gomock.NewController(t)
	accessor := GetUserAuthorizedAccessor(ctrl)
	user := GetTestUser(accessor.database)
	localCard := GetTestLocalCard(accessor.database, user)
	return ctrl, accessor, user, localCard
}

func TestUserAuthorizedAccessorValidAccess(t *testing.T) {
	_, accessor, user, localCard := setupUserAccessorTest(t)

	result, err := accessor.Get(user, &LocalCard{PublicId: localCard.PublicId})
	require.NotNilf(t, result, "accessor returned nil")
	require.Nilf(t, err, "accessor returned non nil error")
	obtainedLocalCard := result.(*LocalCard)
	require.Equal(t, localCard.PublicId, obtainedLocalCard.PublicId, "accessor returned non nil error")
}

func TestUserAuthorizedAccessorInvalidAccess(t *testing.T) {
	_, accessor, user, _ := setupUserAccessorTest(t)
	user2 := GetTestUser(accessor.database)
	localCard2 := GetTestLocalCard(accessor.database, user2)

	result, err := accessor.Get(user, &LocalCard{PublicId: localCard2.PublicId})
	require.Nilf(t, result, "accessor did not return nil")
	require.Equalf(t, ErrNoAccess, err, "accessor returned error other than NoAccess")
}

func TestUserAuthorizedAccessorNoEntity(t *testing.T) {
	_, accessor, user, _ := setupUserAccessorTest(t)

	result, err := accessor.Get(user, &LocalCard{PublicId: "test"})
	require.Nilf(t, result, "accessor did not return nil")
	require.Equalf(t, ErrNotFound, err, "accessor returned error other than NotFound")
}

// TransactionAuthorizedAccessor

func setupAuthorizedTransactionAccessorTest(t *testing.T) (*gomock.Controller, *AuthorizedTransactionAccessorImpl, *User, *User, *Business, *VirtualCard, *Transaction) {
	ctrl := gomock.NewController(t)
	accessor := GetTransactionAuthorizedAccessor(ctrl)
	user := GetTestUser(accessor.database)

	businessUser := GetTestUser(accessor.database)
	business := GetTestBusiness(accessor.database, businessUser)
	itemDefinition := GetTestItemDefinition(accessor.database, business,
		*GetTestFileMetadata(accessor.database, businessUser))

	virtualCard := GetTestVirtualCard(accessor.database, user, business)
	item := GetTestOwnedItem(accessor.database, itemDefinition, virtualCard)
	transaction, _ := GetTestTransaction(accessor.database, virtualCard, []OwnedItem{*item})
	return ctrl, accessor, user, businessUser, business, virtualCard, transaction
}

func TestTransactionAuthorizedAccessorValidFromUser(t *testing.T) {
	_, accessor, user, _, _, _, transaction := setupAuthorizedTransactionAccessorTest(t)

	//println("test virtualcard id ", card.ID, card.OwnerId)
	//println("test user id ", user.ID)
	//println("test transaction virtualcard id ", transaction.VirtualCardId)
	result, err := accessor.GetForUser(user, transaction.Code)
	require.Nilf(t, err, "accessor returned an error")
	require.NotNilf(t, result, "accessor returned nil")
	if result == nil {
		return
	}
	require.Equalf(t, transaction.PublicId, result.PublicId, "accessor returned a different transaction")
	require.Equalf(t, 1, len(result.TransactionDetails), "TransactionDetails number is not 1")
}

func TestTransactionAuthorizedAccessorInvalidFromUser(t *testing.T) {
	_, accessor, _, _, _, _, transaction := setupAuthorizedTransactionAccessorTest(t)
	user2 := GetTestUser(accessor.database)

	result, err := accessor.GetForUser(user2, transaction.Code)
	require.Nilf(t, result, "accessor did not return nil")
	require.Equalf(t, ErrNoAccess, err, "accessor returned error other than NotFound")
}

func TestTransactionAuthorizedAccessorNoEntityFromUser(t *testing.T) {
	_, accessor, user, _, _, _, _ := setupAuthorizedTransactionAccessorTest(t)

	result, err := accessor.GetForUser(user, "asdasd")
	require.Nilf(t, result, "accessor did not return nil")
	require.Equalf(t, ErrNotFound, err, "accessor returned error other than NotFound")
}

func TestTransactionAuthorizedAccessorValidFromBusiness(t *testing.T) {
	_, accessor, _, _, business, _, transaction := setupAuthorizedTransactionAccessorTest(t)

	result, err := accessor.GetForBusiness(business, transaction.Code)
	require.Nilf(t, err, "accessor returned an error")
	require.NotNilf(t, result, "accessor returned nil")
	if result == nil {
		return
	}
	require.Equalf(t, transaction.PublicId, result.PublicId, "accessor returned a different transaction")
	require.Equalf(t, 1, len(result.TransactionDetails), "TransactionDetails number is not 1")
}

func TestTransactionAuthorizedAccessorInvalidFromBusiness(t *testing.T) {
	_, accessor, _, _, _, _, transaction := setupAuthorizedTransactionAccessorTest(t)
	business2 := GetTestBusiness(accessor.database, GetTestUser(accessor.database))

	result, err := accessor.GetForBusiness(business2, transaction.Code)
	require.Nilf(t, result, "accessor returned something other than nil")
	require.Equalf(t, ErrNoAccess, err, "accessor returned a different error than NoAccess")
}

func TestTransactionAuthorizedAccessorNoEntityFromBuiness(t *testing.T) {
	_, accessor, _, _, business, _, _ := setupAuthorizedTransactionAccessorTest(t)

	result, err := accessor.GetForBusiness(business, "123213")
	require.Nilf(t, result, "accessor returned something other than nil")
	require.Equalf(t, ErrNotFound, err, "accessor returned a different error than NotFound")
}
