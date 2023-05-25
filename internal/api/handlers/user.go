package api

import (
	"log"

	"github.com/gin-gonic/gin"

	api "github.com/StampWallet/backend/internal/api/models"
	"github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/database/accessors"
	. "github.com/StampWallet/backend/internal/managers"
	"github.com/StampWallet/backend/internal/services"
	. "github.com/StampWallet/backend/internal/utils"
)

type UserHandlers struct {
	businessManager       BusinessManager
	userAuthorizedAcessor UserAuthorizedAccessor
	virtualCardManager    VirtualCardManager
	localCardManager      LocalCardManager

	// not sure if these two are necessary here
	transactionManager TransactionManager
	localCardHandlers  *UserLocalCardHandlers

	virtualCardHandlers           *UserVirtualCardHandlers
	authorizedTransactionAccessor AuthorizedTransactionAccessor

	logger *log.Logger
}

func CreateUserHandlers(
	virtualCardManager VirtualCardManager,
	localCardManager LocalCardManager,
	businessManager BusinessManager,
	transactionManager TransactionManager,
	itemDefinitionManager ItemDefinitionManager,
	userAuthorizedAcessor UserAuthorizedAccessor,
	authorizedTransactionAccessor AuthorizedTransactionAccessor,
	logger *log.Logger,
) *UserHandlers {
	return &UserHandlers{
		virtualCardManager: virtualCardManager,
		localCardManager:   localCardManager,
		businessManager:    businessManager,
		transactionManager: transactionManager,
		virtualCardHandlers: &UserVirtualCardHandlers{
			virtualCardManager:            virtualCardManager,
			transactionManager:            transactionManager,
			itemDefinitionManager:         itemDefinitionManager,
			userAuthorizedAcessor:         userAuthorizedAcessor,
			authorizedTransactionAccessor: authorizedTransactionAccessor,
			logger:                        services.NewPrefix(logger, "VirtualCardHandlers"),
		},
		localCardHandlers: &UserLocalCardHandlers{
			localCardManager:      localCardManager,
			userAuthorizedAcessor: userAuthorizedAcessor,
			logger:                services.NewPrefix(logger, "UserLocalCardHandlers"),
		},
		userAuthorizedAcessor:         userAuthorizedAcessor,
		authorizedTransactionAccessor: authorizedTransactionAccessor,
		logger:                        logger,
	}
}

func (handler *UserHandlers) getUserCards(c *gin.Context) {

}

func (handler *UserHandlers) getSearchBusinesses(c *gin.Context) {

}

func (handler *UserHandlers) Connect(rg *gin.RouterGroup) {
	cards := rg.Group("/cards")
	{
		handler.localCardHandlers.Connect(cards.Group("/local"))
	}
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

// Handles new local card request
func (handler *UserLocalCardHandlers) postCard(c *gin.Context) {
	req := api.PostUserLocalCardsRequest{}
	if err := c.BindJSON(&req); err != nil {
		handler.logger.Printf("failed to parse in postCard %+v", err)
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST})
		return
	}

	// Get user from context (should be inserted by authMiddleware)
	userAny, exists := c.Get("user")
	if !exists {
		handler.logger.Printf("user not available context")
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}
	user := userAny.(*database.User)

	// Create local card, handle errors
	localCard, err := handler.localCardManager.Create(user, LocalCardDetails{
		Type: req.Type,
		Code: req.Code,
		Name: req.Name,
	})
	if err == ErrInvalidCardType {
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "INVALID_CARD_TYPE"})
		return
	} else if err == ErrCardAlreadyExists {
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "CARD_ALREADY_EXISTS"})
		return
	} else if err != nil {
		handler.logger.Printf("%s unknown error after LocalCardManager.Create: %+v", CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	c.JSON(201, api.PostUserLocalCardsResponse{PublicId: localCard.PublicId})
	return
}

// TODO cache
// Handles get card types request
func (handler *UserLocalCardHandlers) getCardTypes(c *gin.Context) {
	// Convert internal card types to api card types
	var types []api.GetUserLocalCardTypesResponseTypesInner
	for _, v := range CardTypes {
		types = append(types, api.GetUserLocalCardTypesResponseTypesInner{
			PublicId: v.PublicId,
			Name:     v.Name,
			Code:     string(v.Code),
		})
	}

	c.JSON(200, api.GetUserLocalCardTypesResponse{
		Types: types,
	})
}

// Handles delete card request
func (handler *UserLocalCardHandlers) deleteCard(c *gin.Context) {
	cardId := c.Param("cardId")

	// Get user from context (should be inserted by authMiddleware)
	userAny, exists := c.Get("user")
	if !exists {
		handler.logger.Printf("user not available context")
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}
	user := userAny.(*database.User)

	// Get local card form publicId, handle errors
	localCard, err := handler.userAuthorizedAcessor.Get(user, &database.LocalCard{PublicId: cardId})
	if err == ErrNotFound || err == ErrNoAccess {
		c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND})
		return
	} else if err != nil {
		handler.logger.Printf("%s unknown error after userAuthorizedAcessor.Get: %+v", CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	// Remove the card if exists, handle errors
	err = handler.localCardManager.Remove(localCard.(*database.LocalCard))
	if err == ErrCardDoesNotExist {
		c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND})
		return
	} else if err != nil {
		handler.logger.Printf("%s unknown error after localCardManager.Remove: %+v", CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	c.JSON(200, api.DefaultResponse{Status: api.OK})
	return
}

func (handler *UserLocalCardHandlers) Connect(rg *gin.RouterGroup) {
	rg.POST("", handler.postCard)
	rg.GET("/types", handler.getCardTypes)
	rg.DELETE("/:cardId", handler.deleteCard)
}
