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

func getLocalCardManager(ctrl *gomock.Controller) *LocalCardManagerImpl {
    return &LocalCardManagerImpl {
        &BaseServices {
            Logger: log.Default(),
	    Database: GetDatabase(),
        },
    }
}

func TestLocalCardManagerCreate(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager := getLocalCardManager(ctrl)
    user := getTestUser(manager.baseServices.Database)
    localCard, err := manager.Create(user, &LocalCardDetails{
	Type: "s7lJTYHX",
	Code: "012345678901",
    })

    require.Equalf(t, err, nil, "LocalCard.Create returned an error %w", err)
    require.Equalf(t, localCard.Type, "s7lJTYHX", "LocalCard.Create returned wrong card type")
    require.Equalf(t, localCard.Code, "012345678901", "LocalCard.Create returned wrong card code")

    var dbLocalCard LocalCard 
    tx := manager.baseServices.Database.Find(&dbLocalCard, LocalCard{
	Model: gorm.Model { ID: localCard.ID },
    })
    txErr := tx.GetError()
    require.Equalf(t, txErr, nil, "database find returned an error %w", txErr)
    require.Equalf(t, dbLocalCard.Type, "s7lJTYHX", "database has invalid card type")
    require.Equalf(t, dbLocalCard.Code, "012345678901", "database has invalid card number")
}

func TestLocalCardManagerCreateInvalidType(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager := getLocalCardManager(ctrl)
    user := getTestUser(manager.baseServices.Database)
    localCard, err := manager.Create(user, &LocalCardDetails{
	Type: "invalid type lol",
	Code: "012345678901",
    })

    require.Equalf(t, err, InvalidCardType, "LocalCard.Create did not return a InvalidCardType error %w", err)
    require.Equalf(t, localCard, nil, "LocalCard.Create did not return nil LocalCard")
}

func TestLocalCardManagerRemove(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager := getLocalCardManager(ctrl)
    user := getTestUser(manager.baseServices.Database)
    localCard := getTestLocalCard(manager.baseServices.Database, user)

    err := manager.Remove(localCard) 
    require.Equalf(t, err, nil, "LocalCard.Remove returned an error %w", err)

    var dbLocalCard []LocalCard
    tx := manager.baseServices.Database.Find(&dbLocalCard, LocalCard{ Model: gorm.Model { ID: localCard.ID } })
    txErr := tx.GetError()
    require.Equalf(t, txErr, nil, "database find returned an error %w", txErr)
    require.Equalf(t, len(dbLocalCard), 0, "database find returned data")

    nErr := manager.Remove(localCard) 
    require.Equalf(t, nErr, CardDoesNotExist, "LocalCard.Remove did not return a CardDoesNotExist %w", err)
}

func TestLocalCardManagerGetForUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager := getLocalCardManager(ctrl)
    user := getTestUser(manager.baseServices.Database)
    localCard := getTestLocalCard(manager.baseServices.Database, user)

    localCards, err := manager.GetForUser(user) 
    require.Equalf(t, err, nil, "LocalCard.GetForUser returned an error %w", err)
    require.Equalf(t, len(localCards), 1, "LocalCard.GetForUser returned more or less than 1 %d", 
	len(localCards))
    require.Equalf(t, localCards[0].PublicId, localCard.PublicId, 
	"LocalCard.GetForUser returned unexpected card %s != %s", 
	localCards[0].PublicId, localCard.PublicId)
}
