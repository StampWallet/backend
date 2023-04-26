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

func getTransactionManager(ctrl *gomock.Controller) *TransactionManagerImpl {
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
    //manager := getTransactionManager(ctrl)
    //user := getTestUser(manager.baseServices.Database)
    //businessUser := getTestUser(manager.baseServices.Database)
    //business := getTestBusiness(manager.baseServices.Database, businessUser)
    //virtualCard := getTestVirtualCard(manager.baseServices.Database, user, business)
}

func TestTransactionManagerFinalize(t *testing.T) {
    //ctrl := gomock.NewController(t)
    //defer ctrl.Finish()
    //manager := getTransactionManager(ctrl)
}
