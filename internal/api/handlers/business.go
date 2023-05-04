package api

import (
	"log"

	. "github.com/StampWallet/backend/internal/database/accessors"
	. "github.com/StampWallet/backend/internal/managers"
	"github.com/gin-gonic/gin"
)

type BusinessHandlers struct {
	businessManager               BusinessManager
	transactionManager            TransactionManager
	itemDefinitionHandlers        ItemDefinitionHandlers
	userAuthorizedAcessor         UserAuthorizedAccessor
	businessAuthorizedAccessor    BusinessAuthorizedAccessor
	authorizedTransactionAccessor AuthorizedTransactionAccessor
	logger                        *log.Logger
}

func (handler *BusinessHandlers) postAccountInfo(c *gin.Context) {

}

func (handler *BusinessHandlers) getAccountInfo(c *gin.Context) {

}

func (handler *BusinessHandlers) patchAccountInfo(c *gin.Context) {

}

func (handler *BusinessHandlers) getTransaction(c *gin.Context) {

}

func (handler *BusinessHandlers) postTransaction(c *gin.Context) {

}

func (handler *BusinessHandlers) Connect(rg *gin.RouterGroup) {

}

type ItemDefinitionHandlers struct {
	itemDefinitionManager      ItemDefinitionManager
	userAuthorizedAcessor      UserAuthorizedAccessor
	businessAuthorizedAccessor BusinessAuthorizedAccessor
	logger                     *log.Logger
}

func (handler *ItemDefinitionHandlers) getItemDefinition(c *gin.Context) {

}

func (handler *ItemDefinitionHandlers) postItemDefinition(c *gin.Context) {

}

func (handler *ItemDefinitionHandlers) patchItemDefinition(c *gin.Context) {

}

func (handler *ItemDefinitionHandlers) deleteItemDefinition(c *gin.Context) {

}

func (handler *ItemDefinitionHandlers) Connect(rg *gin.RouterGroup) {

}
