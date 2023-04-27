package managers

import (
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services"
	. "github.com/StampWallet/backend/internal/testutils"
)

func GetTransactionManager(ctrl *gomock.Controller) *TransactionManagerImpl {
	return &TransactionManagerImpl{
		&BaseServices{
			Logger:   log.Default(),
			Database: GetDatabase(),
		},
	}
}

type transactionTest struct {
	ctrl           *gomock.Controller
	manager        *TransactionManagerImpl
	user           *User
	businessUser   *User
	business       *Business
	itemDefinition *ItemDefinition
	virtualCard    *VirtualCard
	ownedItem      *OwnedItem
	db             GormDB
}

func setupTransactionTest(t *testing.T) transactionTest {
	ctrl := gomock.NewController(t)
	manager := GetTransactionManager(ctrl)
	db := manager.baseServices.Database
	user := GetTestUser(db)
	businessUser := GetTestUser(db)
	business := GetTestBusiness(db, businessUser)
	itemDefinition := GetTestItemDefinition(db, business, *GetTestFileMetadata(db, user))
	virtualCard := GetTestVirtualCard(db, user, business)
	ownedItem := GetTestOwnedItem(db, itemDefinition, virtualCard)
	return transactionTest{
		ctrl:           ctrl,
		manager:        manager,
		user:           user,
		businessUser:   businessUser,
		business:       business,
		itemDefinition: itemDefinition,
		virtualCard:    virtualCard,
		ownedItem:      ownedItem,
		db:             db,
	}
}

func TestTransactionManagerStart(t *testing.T) {
	s := setupTransactionTest(t)

	transaction, err := s.manager.Start(s.virtualCard, []OwnedItem{*s.ownedItem})
	require.Nilf(t, err, "transaction start returned an error %w", err)
	if transaction == nil {
		t.Errorf("transacition is nil")
		return
	}
	require.Equalf(t, TransactionStateStarted, transaction.State, "TransactionManager.Start returned transaction with invalid state %s", transaction.State)
	require.NotEqualf(t, 0, transaction.AddedPoints, "TransactionManager.Start returned transaction with more than 0 points")
	require.NotNilf(t, transaction.Code, "TransactionManager.Start returned transaction with nil code")

	var dbTransaction []Transaction
	tx := s.db.Find(&transaction, Transaction{Model: gorm.Model{ID: transaction.ID}})
	err = tx.GetError()
	require.Nilf(t, err, "transaction find returned an error %w", err)
	require.Nilf(t, dbTransaction[0], "TransactionManager.Start returned nil transacition")
	require.Equalf(t, TransactionStateStarted, dbTransaction[0].State, "TransactionManager.Start returned transaction with invalid state %s", dbTransaction[0].State)
	require.NotEqualf(t, 0, dbTransaction[0].AddedPoints, "TransactionManager.Start returned transaction with more than 0 points")
	require.NotNilf(t, dbTransaction[0].Code, "TransactionManager.Start returned transaction with nil code")

	var transactionDetails []TransactionDetail
	tx = s.db.Find(&transactionDetails, TransactionDetail{TransactionId: transaction.ID})
	err = tx.GetError()
	require.Nilf(t, err, "database find for TransactionDetails returned an error %w", err)
	require.Equalf(t, 1, len(transactionDetails),
		"database find for TransactionDetails returned less or more than 1 row %d", len(transactionDetails))
	require.Equalf(t, s.ownedItem.ID, transactionDetails[0].ItemId,
		"database find for TransactionDetails returned an invalid item %d != %d",
		s.ownedItem.ID, transactionDetails[0].ItemId)
	require.Equalf(t, NoActionType, transactionDetails[0].Action,
		"database find for TransactionDetails returned invalid action for itmem %s",
		transactionDetails[0].Action)
}

func TestTransactionManagerStartWithInvalidItems(t *testing.T) {
	s := setupTransactionTest(t)
	invalidStatus := []OwnedItemStatusEnum{OwnedItemStatusUsed, OwnedItemStatusReturned, OwnedItemStatusWithdrawn}
	for _, status := range invalidStatus {
		ownedItem := GetDefaultOwnedItem(s.itemDefinition, s.virtualCard)
		ownedItem.Status = status
		Save(s.db, ownedItem)
		transaction, err := s.manager.Start(s.virtualCard, []OwnedItem{*s.ownedItem})
		require.Nilf(t, transaction, "TransactionManager.Start should return a nil transaction")
		require.Equalf(t, err, InvalidItem, "TransactionManager.Start should return a InvalidItem error")
	}
}

//TODO transaction expiration? status exists

func TestTransactionManagerFinalize(t *testing.T) {
	s := setupTransactionTest(t)
	ownedItemToRedeem := GetTestOwnedItem(s.db, s.itemDefinition, s.virtualCard)
	ownedItemToRecall := GetTestOwnedItem(s.db, s.itemDefinition, s.virtualCard)
	ownedItemToCancel := GetTestOwnedItem(s.db, s.itemDefinition, s.virtualCard)
	transaction, _ := GetTestTransaction(s.db, s.virtualCard,
		[]OwnedItem{*ownedItemToRedeem, *ownedItemToRecall, *ownedItemToCancel})

	transaction, err := s.manager.Finalize(transaction, []ItemWithAction{
		{ownedItemToRedeem, RedeemedActionType},
		{ownedItemToRecall, RecalledActionType},
		{ownedItemToCancel, CancelledActionType},
	}, 10)
	require.Nilf(t, err, "transaction finalize returned an error %w", err)
	if transaction == nil {
		t.Errorf("transaction is nil")
		return
	}
	require.Equalf(t, 10, transaction.AddedPoints, "transaction has a different number of added points. Expected: %d, got %d", 10, transaction.AddedPoints)
	require.Equalf(t, TransactionStateFinished, transaction.State, "transaction has a different state (%s) than expected (%s)", transaction.State, TransactionStateFinished)

	var dbTransaction Transaction
	tx := s.db.First(&dbTransaction, Transaction{Model: gorm.Model{ID: transaction.ID}})
	err = tx.GetError()
	require.Nilf(t, err, "database find for Transaction returned an error %w", err)
	require.Equalf(t, 10, dbTransaction.AddedPoints, "db transaction has a different number of added points. Expected: %d, got %d", 10, dbTransaction.AddedPoints)
	require.Equalf(t, TransactionStateFinished, dbTransaction.State, "transaction has a different state (%s) than expected (%s)", dbTransaction.State, TransactionStateFinished)

	var dbTransactionDetails []TransactionDetail
	tx = s.db.Find(&dbTransaction, TransactionDetail{TransactionId: transaction.ID})
	err = tx.GetError()
	require.Nilf(t, err, "database find for TransactionDetails returned an error %w", err)
	require.Equalf(t, 3, len(dbTransactionDetails), "db returned a different number of transaction details than expected. Expected: %d, got %d", 3, len(dbTransactionDetails))
	for _, d := range dbTransactionDetails {
		var dbOwnedItem OwnedItem
		tx := s.db.Find(&dbOwnedItem, OwnedItem{Model: gorm.Model{ID: d.ItemId}})
		err := tx.GetError()
		require.Nilf(t, err, "database find for OwnedItem returned an error %w", err)

		var expectedActionType ActionTypeEnum
		switch d.ItemId {
		case ownedItemToRedeem.ID:
			expectedActionType = RedeemedActionType
		case ownedItemToRecall.ID:
			expectedActionType = RecalledActionType
		case ownedItemToCancel.ID:
			expectedActionType = CancelledActionType
		default:
			t.Errorf("unexpected transaction detail %d", d.ID)
			continue
		}

		require.Equalf(t, expectedActionType, d.Action, "invalid action type %s != %s",
			d.Action, expectedActionType)
		if expectedActionType == RedeemedActionType {
			require.Truef(t, dbOwnedItem.Used.Valid, "owned item used time not valid")
			require.Equalf(t, OwnedItemStatusUsed,
				dbOwnedItem.Status, "owned item status does not equal OwnedItemStatusUsed but %s", dbOwnedItem.Status)
		} else if expectedActionType == RecalledActionType {
			require.Truef(t, dbOwnedItem.Used.Valid, "owned item used time not valid")
			require.Equalf(t, OwnedItemStatusWithdrawn,
				dbOwnedItem.Status, "owned item status does not equal OwnedItemStatusWithdrawn but %s", dbOwnedItem.Status)
		}
	}

	var dbVirtualCard VirtualCard
	tx = s.db.First(&dbVirtualCard, VirtualCard{Model: gorm.Model{ID: s.virtualCard.ID}})
	err = tx.GetError()
	require.Nilf(t, err, "database find for TransactionDetails returned an error %w", err)
	require.Equalf(t, s.virtualCard.Points+s.itemDefinition.Price+10, dbVirtualCard.Points,
		"virtual card has a wrong number of points. Expected: %d Got: %d",
		s.virtualCard.Points+s.itemDefinition.Price+10, dbVirtualCard.Points)
}

func TestTransactionManagerFinalizeWithItemsNotFromTransaction(t *testing.T) {
	s := setupTransactionTest(t)
	ownedItemToRedeem := GetTestOwnedItem(s.db, s.itemDefinition, s.virtualCard)
	ownedItemFromOutside := GetTestOwnedItem(s.db, s.itemDefinition, s.virtualCard)
	transaction, _ := GetTestTransaction(s.db, s.virtualCard,
		[]OwnedItem{*ownedItemToRedeem})

	transaction, err := s.manager.Finalize(transaction, []ItemWithAction{
		{ownedItemFromOutside, RedeemedActionType},
	}, 10)
	require.Equalf(t, InvalidItem, err, "TransactionManager.Finalize did not return InvalidItemError %w",
		InvalidItem)

	var dbTransaction Transaction
	tx := s.db.First(&dbTransaction, Transaction{Model: gorm.Model{ID: transaction.ID}})
	err = tx.GetError()
	require.Nilf(t, err, "database find for Transaction returned an error %w", err)
	require.Equalf(t, TransactionStateStarted, dbTransaction.State, "dbTransaction state is not TransactionStateStarted %s", dbTransaction.State)

	var dbTransactionDetail TransactionDetail
	tx = s.db.First(&dbTransactionDetail, TransactionDetail{TransactionId: transaction.ID})
	err = tx.GetError()
	require.Nilf(t, err, "database find for TransactionDetail returned an error %w", err)
	require.Equalf(t, NoActionType, dbTransactionDetail.Action,
		"dbTransactionDetail action is not NoActionType %s", dbTransactionDetail.Action)
}
