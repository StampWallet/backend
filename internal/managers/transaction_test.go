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
    return &TransactionManagerImpl {
        &BaseServices {
            Logger: log.Default(),
	    Database: GetDatabase(),
        },
    }
}

func TestTransactionManagerStart(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager := GetTransactionManager(ctrl)
    db := manager.baseServices.Database
    user := GetTestUser(db)
    businessUser := GetTestUser(db)
    business := GetTestBusiness(db, businessUser)
    itemDefinition := GetTestItemDefinition(db, business, *GetTestFileMetadata(db, user))
    virtualCard := GetTestVirtualCard(db, user, business)
    ownedItem := GetTestOwnedItem(db, itemDefinition, virtualCard)

    transaction, err := manager.Start(virtualCard, []OwnedItem{ *ownedItem }) 
    require.Equalf(t, err, nil, "transaction start returned an error %w", err)
    require.Equalf(t, transaction, nil, "TransactionManager.Start returned nil transacition")
    require.Equalf(t, transaction.State, TransactionStateStarted, "TransactionManager.Start returned transaction with invalid state %s", transaction.State)
    require.NotEqualf(t, transaction.AddedPoints, 0, "TransactionManager.Start returned transaction with more than 0 points")
    require.NotEqualf(t, transaction.Code, nil, "TransactionManager.Start returned transaction with nil code")

    var dbTransaction []Transaction
    tx := db.Find(&transaction, Transaction{ Model: gorm.Model { ID: transaction.ID } })
    err = tx.GetError()
    require.Equalf(t, err, nil, "transaction find returned an error %w", err)
    require.Equalf(t, dbTransaction[0], nil, "TransactionManager.Start returned nil transacition")
    require.Equalf(t, dbTransaction[0].State, TransactionStateStarted, "TransactionManager.Start returned transaction with invalid state %s", dbTransaction[0].State)
    require.NotEqualf(t, dbTransaction[0].AddedPoints, 0, "TransactionManager.Start returned transaction with more than 0 points")
    require.NotEqualf(t, dbTransaction[0].Code, nil, "TransactionManager.Start returned transaction with nil code")

    var transactionDetails []TransactionDetail
    tx = db.Find(&transactionDetails, TransactionDetail{ TransactionId: transaction.ID })
    err = tx.GetError()
    require.Equalf(t, err, nil, "database find for TransactionDetails returned an error %w", err)
    require.Equalf(t, len(transactionDetails), 1, 
	"database find for TransactionDetails returned less or more than 1 row %d", len(transactionDetails))
    require.Equalf(t, transactionDetails[0].ItemId, ownedItem.ID, 
	"database find for TransactionDetails returned an invalid item %d != %d", 
	ownedItem.ID, transactionDetails[0].ItemId)
    require.Equalf(t, transactionDetails[0].Action, NoActionType, 
	"database find for TransactionDetails returned invalid action for itmem %s", 
	transactionDetails[0].Action)
}

//TODO transaction expiration? status exists

func TestTransactionManagerFinalize(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager := GetTransactionManager(ctrl)
    db := manager.baseServices.Database
    user := GetTestUser(db)
    businessUser := GetTestUser(db)
    business := GetTestBusiness(db, businessUser)
    itemDefinition := GetTestItemDefinition(db, business, *GetTestFileMetadata(db, user))
    virtualCard := GetTestVirtualCard(db, user, business)
    ownedItemToRedeem := GetTestOwnedItem(db, itemDefinition, virtualCard)
    ownedItemToRecall := GetTestOwnedItem(db, itemDefinition, virtualCard)
    ownedItemToCancel := GetTestOwnedItem(db, itemDefinition, virtualCard)
    transaction, _ := GetTestTransaction(db, virtualCard, 
	[]OwnedItem{ *ownedItemToRedeem, *ownedItemToRecall, *ownedItemToCancel })

    transaction, err := manager.Finalize(transaction, []ItemWithAction{ 
	{ ownedItemToRedeem, RedeemedActionType },
	{ ownedItemToRecall, RecalledActionType },
	{ ownedItemToCancel, CancelledActionType },
    }, 10)
    require.Equalf(t, err, nil, "transaction finalize returned an error %w", err)
    require.Equalf(t, transaction.AddedPoints, 10, "transaction has a different number of added points. Expected: %d, got %d", 10, transaction.AddedPoints)
    require.Equalf(t, transaction.State, TransactionStateFinished, "transaction has a different state (%s) than expected (%s)", transaction.State, TransactionStateFinished)

    var dbTransaction Transaction
    tx := db.First(&dbTransaction, Transaction{ Model: gorm.Model{ ID: transaction.ID }})
    err = tx.GetError()
    require.Equalf(t, err, nil, "database find for Transaction returned an error %w", err)
    require.Equalf(t, dbTransaction.AddedPoints, 10, "db transaction has a different number of added points. Expected: %d, got %d", 10, dbTransaction.AddedPoints)
    require.Equalf(t, dbTransaction.State, TransactionStateFinished, "transaction has a different state (%s) than expected (%s)", dbTransaction.State, TransactionStateFinished)

    var dbTransactionDetails []TransactionDetail
    tx = db.Find(&dbTransaction, TransactionDetail{ TransactionId: transaction.ID })
    err = tx.GetError()
    require.Equalf(t, err, nil, "database find for TransactionDetails returned an error %w", err)
    require.Equalf(t, len(dbTransactionDetails), 3, "db returned a different number of transaction details than expected. Expected: %d, got %d", 3, len(dbTransactionDetails))
    for _, d := range dbTransactionDetails {
	var dbOwnedItem OwnedItem
	tx := db.Find(&dbOwnedItem, OwnedItem{ Model: gorm.Model { ID: d.ItemId } })
	err := tx.GetError()
	require.Equalf(t, err, nil, "database find for OwnedItem returned an error %w", err)

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

	require.Equalf(t, d.Action, expectedActionType, "invalid action type %s != %s", 
	    d.Action, expectedActionType)
	if expectedActionType == RedeemedActionType {
	    require.Truef(t, dbOwnedItem.Used.Valid, "owned item used time not valid")
	    require.Equalf(t, dbOwnedItem.Status, OwnedItemStatusUsed, 
		"owned item status does not equal OwnedItemStatusUsed but %s", dbOwnedItem.Status)
	} else if expectedActionType == RecalledActionType {
	    require.Truef(t, dbOwnedItem.Used.Valid, "owned item used time not valid")
	    require.Equalf(t, dbOwnedItem.Status, OwnedItemStatusWithdrawn, 
		"owned item status does not equal OwnedItemStatusWithdrawn but %s", dbOwnedItem.Status)
	}
    }

    var dbVirtualCard VirtualCard
    tx = db.First(&dbVirtualCard, VirtualCard{ Model: gorm.Model{ ID: virtualCard.ID }})
    err = tx.GetError()
    require.Equalf(t, err, nil, "database find for TransactionDetails returned an error %w", err)
    require.Equalf(t, dbVirtualCard.Points, virtualCard.Points + itemDefinition.Price + 10, 
	"virtual card has a wrong number of points. Expected: %d Got: %d", 
	virtualCard.Points + itemDefinition.Price + 10, dbVirtualCard.Points)
}

func TestTransactionManagerFinalizeWithItemsNotFromTransaction(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager := GetTransactionManager(ctrl)
    db := manager.baseServices.Database
    user := GetTestUser(db)
    businessUser := GetTestUser(db)
    business := GetTestBusiness(db, businessUser)
    itemDefinition := GetTestItemDefinition(db, business, *GetTestFileMetadata(db, user))
    virtualCard := GetTestVirtualCard(db, user, business)
    ownedItemToRedeem := GetTestOwnedItem(db, itemDefinition, virtualCard)
    ownedItemFromOutside := GetTestOwnedItem(db, itemDefinition, virtualCard)
    transaction, _ := GetTestTransaction(db, virtualCard, 
	[]OwnedItem{ *ownedItemToRedeem })

    transaction, err := manager.Finalize(transaction, []ItemWithAction{ 
	{ ownedItemFromOutside, RedeemedActionType },
    }, 10)
    require.Equalf(t, err, InvalidItem, "TransactionManager.Finalize did not return InvalidItemError %w", 
	InvalidItem)

    var dbTransaction Transaction
    tx := db.First(&dbTransaction, Transaction{ Model: gorm.Model{ ID: transaction.ID }})
    err = tx.GetError()
    require.Equalf(t, err, nil, "database find for Transaction returned an error %w", err)
    require.Equalf(t, dbTransaction.State, TransactionStateStarted, "dbTransaction state is not TransactionStateStarted %s", dbTransaction.State)

    var dbTransactionDetail TransactionDetail
    tx = db.First(&dbTransactionDetail, TransactionDetail{ TransactionId: transaction.ID })
    err = tx.GetError()
    require.Equalf(t, err, nil, "database find for TransactionDetail returned an error %w", err)
    require.Equalf(t, dbTransactionDetail.Action, NoActionType, 
	"dbTransactionDetail action is not NoActionType %s", dbTransactionDetail.Action)
}
