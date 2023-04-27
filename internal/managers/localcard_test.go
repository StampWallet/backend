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

func GetLocalCardManager(ctrl *gomock.Controller) *LocalCardManagerImpl {
	return &LocalCardManagerImpl{
		&BaseServices{
			Logger:   log.Default(),
			Database: GetDatabase(),
		},
	}
}

func TestLocalCardManagerCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetLocalCardManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	localCard, err := manager.Create(user, &LocalCardDetails{
		Type: "s7lJTYHX",
		Code: "012345678901",
	})

	require.Nilf(t, err, "LocalCard.Create returned an error %w", err)
	if localCard == nil {
		t.Errorf("local card is nil")
		return
	}
	require.Equalf(t, "s7lJTYHX", localCard.Type, "LocalCard.Create returned wrong card type")
	require.Equalf(t, "012345678901", localCard.Code, "LocalCard.Create returned wrong card code")

	var dbLocalCard LocalCard
	tx := manager.baseServices.Database.Find(&dbLocalCard, LocalCard{
		Model: gorm.Model{ID: localCard.ID},
	})
	txErr := tx.GetError()
	require.Nilf(t, txErr, "database find returned an error %w", txErr)
	require.Equalf(t, "s7lJTYHX", dbLocalCard.Type, "database has invalid card type")
	require.Equalf(t, "012345678901", dbLocalCard.Code, "database has invalid card number")
}

func TestLocalCardManagerCreateInvalidType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetLocalCardManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	localCard, err := manager.Create(user, &LocalCardDetails{
		Type: "invalid type lol",
		Code: "012345678901",
	})

	require.Equalf(t, InvalidCardType, err, "LocalCard.Create did not return a InvalidCardType error %w", err)
	require.Nilf(t, localCard, "LocalCard.Create did not return nil LocalCard")
}

func TestLocalCardManagerRemove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetLocalCardManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	localCard := GetTestLocalCard(manager.baseServices.Database, user)

	err := manager.Remove(localCard)
	require.Nilf(t, err, "LocalCard.Remove returned an error %w", err)

	var dbLocalCard []LocalCard
	tx := manager.baseServices.Database.Find(&dbLocalCard, LocalCard{Model: gorm.Model{ID: localCard.ID}})
	txErr := tx.GetError()
	require.Nilf(t, txErr, "database find returned an error %w", txErr)
	require.Equalf(t, 0, len(dbLocalCard), "database find returned data when no data was expected")

	nErr := manager.Remove(localCard)
	require.Equalf(t, CardDoesNotExist, nErr, "LocalCard.Remove did not return a CardDoesNotExist %w", err)
}

func TestLocalCardManagerGetForUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager := GetLocalCardManager(ctrl)
	user := GetTestUser(manager.baseServices.Database)
	localCard := GetTestLocalCard(manager.baseServices.Database, user)

	localCards, err := manager.GetForUser(user)
	require.Nilf(t, err, "LocalCard.GetForUser returned an error %w", err)
	require.Equalf(t, 1, len(localCards), "LocalCard.GetForUser returned more or less than 1 %d",
		len(localCards))
	require.Equalf(t, localCard.PublicId, localCards[0].PublicId,
		"LocalCard.GetForUser returned unexpected card %s != %s",
		localCards[0].PublicId, localCard.PublicId)
}
