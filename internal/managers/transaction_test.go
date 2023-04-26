package managers

import (
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	//"github.com/stretchr/testify/require"
	//"gorm.io/gorm"

	//. "github.com/StampWallet/backend/internal/database"
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
    //ctrl := gomock.NewController(t)
    //defer ctrl.Finish()
    //manager := GetTransactionManager(ctrl)
    //db := manager.baseServices.Database
    //user := GetTestUser(db)
    //businessUser := GetTestUser(db)
    //business := GetTestBusiness(db, businessUser)
    //itemDefinition := GetTestItemDefinition(db, business, *GetTestFileMetadata(db, user))
    //virtualCard := GetTestVirtualCard(db, user, business)
    //ownedItem := GetTestOwnedItem(db, itemDefinition, virtualCard)
}

func TestTransactionManagerFinalize(t *testing.T) {
    //ctrl := gomock.NewController(t)
    //defer ctrl.Finish()
    //manager := GetTransactionManager(ctrl)
}
