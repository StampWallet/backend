package api

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	api "github.com/StampWallet/backend/internal/api/models"
	"github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/database/accessors"
	. "github.com/StampWallet/backend/internal/managers"
	"github.com/StampWallet/backend/internal/services"
	. "github.com/StampWallet/backend/internal/utils"
)

// UserHandlers stores UserLocalCardHandlers and UserVirtualCardHandlers. It also implements
// few requests related to functionalities for users that did not fit other handlers.
// (currently retrieving a list of local and virtual cards, searching for businesses)
// Retrieval of local and virtual cards is coupled in a single request to limit number
// of requests made by the application.
// All requests require middleware that will insert user object into the context under "user".
// Such middleware has to be set up by owner of the external router group
// (whoever calls UserHandlers.Connect).
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

// Creates UserHandlers. UserHandlers "owns" UserVirtualCardHandlers and UserLocalCardHandlers,
// hence these two structs are created in this function, not passed as arguments.
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

// Converts database.Business to api.ShortBusinessDetailsApiModel
// Most data is lost in conversion - api.ShortBusinessDetailsApiModel does not contain all
// data from model
func convertBusinessToShortApiModel(business *database.Business) api.ShortBusinessDetailsApiModel {
	return api.ShortBusinessDetailsApiModel{
		PublicId:       business.PublicId,
		Name:           business.Name,
		Description:    business.Description,
		GpsCoordinates: business.GPSCoordinates.ToString(),
		BannerImageId:  business.BannerImageId,
		IconImageId:    business.IconImageId,
	}
}

// Handles local and virtual cards retrieval request
func (handler *UserHandlers) getUserCards(c *gin.Context) {
	// Get user object from middleware
	userAny, exists := c.Get("user")
	if !exists {
		handler.logger.Printf("user not available context")
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}
	user := userAny.(*database.User)

	result := api.GetUserCardsResponse{}

	localCards, err := handler.userAuthorizedAcessor.GetAll(user, &database.LocalCard{}, []string{})
	if err != nil {
		handler.logger.Printf("%s unknown error after userAuthorizedAcessor.GetAll for localCard: %+v",
			CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	for _, v := range localCards {
		card := v.(*database.LocalCard)
		result.LocalCards = append(result.LocalCards, api.LocalCardApiModel{
			PublicId: card.PublicId,
			Name:     card.Name,
			Type:     card.Type,
			Code:     card.Code,
		})
	}

	virtualCards, err := handler.userAuthorizedAcessor.GetAll(user, &database.VirtualCard{},
		[]string{"Business"})
	if err != nil {
		handler.logger.Printf("%s unknown error after userAuthorizedAcessor.GetAll for virtualCard: %+v",
			CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	for _, v := range virtualCards {
		card := v.(*database.VirtualCard)
		result.VirtualCards = append(result.VirtualCards, api.ShortVirtualCardApiModel{
			BusinessDetails: convertBusinessToShortApiModel(card.Business),
			Points:          int32(card.Points),
		})
	}

	c.JSON(200, result)
}

// Handles business search request
// Requires middleware that will insert user object into the context under "user".
func (handler *UserHandlers) getSearchBusinesses(c *gin.Context) {
	textQuery := c.Query("text")
	locationQuery := c.Query("location")
	proximityQuery, proximityExists := c.GetQuery("proximity")
	offsetQuery := c.Query("offset")
	limitQuery := c.Query("limit")

	var text *string
	if textQuery != "" {
		text = &textQuery
	}

	var location *database.GPSCoordinates
	if locationQuery != "" {
		coords, err := database.GPSCoordinatesFromString(locationQuery)
		if err != nil {
			c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "INVALID_LOCATION"})
			return
		}
		location = &coords
	}

	var proximity uint = 50
	if location != nil && !proximityExists {
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "LOCATION_BUT_NO_PROXIMITY"})
		return
	} else if location != nil {
		localProximity, err := strconv.ParseUint(proximityQuery, 10, 32)
		if err != nil {
			c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "INVALID_PROXIMITY"})
			return
		}
		proximity = uint(localProximity)
	}

	var offset uint = 0
	var limit uint = 50

	if offsetQuery != "" {
		localOffset, err := strconv.ParseUint(offsetQuery, 10, 32)
		if err != nil {
			c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "INVALID_OFFSET"})
			return
		}
		offset = uint(localOffset)
	}

	if limitQuery != "" {
		localLimit, err := strconv.ParseUint(limitQuery, 10, 32)
		if err != nil {
			c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "INVALID_OFFSET"})
			return
		}
		offset = uint(localLimit)
	}

	businesses, err := handler.businessManager.Search(text, location, proximity, offset, limit)
	if err != nil {
		handler.logger.Printf("%s unknown error after businessManager.Search %+v", CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	result := api.GetUserBusinessesSearchResponse{}

	for _, v := range businesses {
		result.Businesses = append(result.Businesses, convertBusinessToShortApiModel(&v))
	}

	c.JSON(200, result)
	return
}

// Connects handler to gin router
func (handler *UserHandlers) Connect(rg *gin.RouterGroup) {
	cards := rg.Group("/cards")
	{
		cards.GET("", handler.getUserCards)
		handler.localCardHandlers.Connect(cards.Group("/local"))
	}
	rg.GET("/businesses", handler.getSearchBusinesses)
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

// UserLocalCardHandlers implements API handlers for requests related to
// managing local cards (creating, deleting, retrieving types of cards).
// All requests require middleware that will insert user object into the context under "user".
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

// Connects UserLocalCardHandlers to gin router
func (handler *UserLocalCardHandlers) Connect(rg *gin.RouterGroup) {
	rg.POST("", handler.postCard)
	rg.GET("/types", handler.getCardTypes)
	rg.DELETE("/:cardId", handler.deleteCard)
}
