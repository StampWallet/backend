package api

import (
	"log"

	. "github.com/StampWallet/backend/internal/database/accessors"
	. "github.com/StampWallet/backend/internal/managers"
	"github.com/gin-gonic/gin"
)

type UserHandlers struct {
	virtualCardManager            VirtualCardManager
	localCardManager              LocalCardManager
	businessManager               BusinessManager
	transactionManager            TransactionManager
	virtualCardHandlers           *UserVirtualCardHandlers
	localCardHandlers             *UserLocalCardHandlers
	userAuthorizedAcessor         UserAuthorizedAccessor
	authorizedTransactionAccessor AuthorizedTransactionAccessor

	logger *log.Logger
}

func (handler *UserHandlers) getUserCards(c *gin.Context) {

}

func (handler *UserHandlers) getSearchBusinesses(c *gin.Context) {

}

func (handler *UserHandlers) Connect(rg *gin.RouterGroup) {

}

type UserVirtualCardHandlers struct {
	virtualCardManager            VirtualCardManager
	transactionManager            TransactionManager
	itemDefinitionManager         ItemDefinitionManager
	userAuthorizedAcessor         UserAuthorizedAccessor
	authorizedTransactionAccessor AuthorizedTransactionAccessor
	logger                        *log.Logger
}

func (handler *UserVirtualCardHandlers) postCard(c *gin.Context) {

}

func (handler *UserVirtualCardHandlers) deleteCard(c *gin.Context) {

}

func (handler *UserVirtualCardHandlers) getCard(c *gin.Context) {

}

func (handler *UserVirtualCardHandlers) postItem(c *gin.Context) {

}

func (handler *UserVirtualCardHandlers) deleteItem(c *gin.Context) {

}

func (handler *UserVirtualCardHandlers) postTransaction(c *gin.Context) {

}

func (handler *UserVirtualCardHandlers) Connect(rg *gin.RouterGroup) {

}

type UserLocalCardHandlers struct {
	localCardManager      LocalCardManager
	userAuthorizedAcessor UserAuthorizedAccessor
	logger                *log.Logger
}

func (handler *UserLocalCardHandlers) postCard(c *gin.Context) {

}

func (handler *UserLocalCardHandlers) getCardTypes(c *gin.Context) {

}

func (handler *UserLocalCardHandlers) deleteCard(c *gin.Context) {

}

func (handler *UserLocalCardHandlers) Connect(rg *gin.RouterGroup) {

}
