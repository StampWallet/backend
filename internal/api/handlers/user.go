package api

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	api "github.com/StampWallet/backend/internal/api/models"
	apiUtils "github.com/StampWallet/backend/internal/api/utils"
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

// Handles local and virtual cards retrieval request
func (handler *UserHandlers) getUserCards(c *gin.Context) {
	// Get user object from middleware
	user := getUserFromContext(handler.logger, c)
	if user == nil {
		return
	}

	result := api.GetUserCardsResponse{}

	// Get all local cards of user
	localCards, err := handler.userAuthorizedAcessor.GetAll(user, &database.LocalCard{}, []string{})
	if err != nil {
		handler.logger.Printf("%s unknown error after userAuthorizedAcessor.GetAll for localCard: %+v",
			CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	// Convert database.LocalCard to api model
	for _, v := range localCards {
		card := v.(*database.LocalCard)
		result.LocalCards = append(result.LocalCards, api.LocalCardApiModel{
			PublicId: card.PublicId,
			Name:     card.Name,
			Type:     card.Type,
			Code:     card.Code,
		})
	}

	// Get all virtual cards of user
	virtualCards, err := handler.userAuthorizedAcessor.GetAll(user, &database.VirtualCard{},
		[]string{"Business"})
	if err != nil {
		handler.logger.Printf("%s unknown error after userAuthorizedAcessor.GetAll for virtualCard: %+v",
			CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	// Convert database.VirtualCard to api model
	for _, v := range virtualCards {
		card := v.(*database.VirtualCard)
		result.VirtualCards = append(result.VirtualCards, api.ShortVirtualCardApiModel{
			BusinessDetails: apiUtils.ConvertBusinessToShortApiModel(card.Business),
			Points:          int32(card.Points),
		})
	}

	c.JSON(200, result)
}

// Handles business search request
func (handler *UserHandlers) getSearchBusinesses(c *gin.Context) {
	// Text search
	textQuery := c.Query("text")
	// Filter by location (long,lat)
	locationQuery := c.Query("location")
	// Filter by location - proximity in meters
	proximityQuery, proximityExists := c.GetQuery("proximity")
	// Pagination - offset
	offsetQuery := c.Query("offset")
	// Pagination - limit
	limitQuery := c.Query("limit")

	// Parse text query if presetn
	var text *string
	if textQuery != "" {
		text = &textQuery
	}

	// Parse location if present
	var location *database.GPSCoordinates
	if locationQuery != "" {
		coords, err := database.GPSCoordinatesFromString(locationQuery)
		if err != nil {
			c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "INVALID_LOCATION"})
			return
		}
		location = &coords
	}

	// Default proximity - 1000 meters
	var proximity uint = 1000
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

	// Parse offset if exists
	if offsetQuery != "" {
		localOffset, err := strconv.ParseUint(offsetQuery, 10, 32)
		if err != nil {
			c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "INVALID_OFFSET"})
			return
		}
		offset = uint(localOffset)
	}

	// Parse limit if exists
	if limitQuery != "" {
		localLimit, err := strconv.ParseUint(limitQuery, 10, 32)
		if err != nil {
			c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "INVALID_OFFSET"})
			return
		}
		offset = uint(localLimit)
	}

	if location == nil && text == nil {
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "EMPTY_QUERY"})
		return
	}

	// Execute the query, handle errors
	businesses, err := handler.businessManager.Search(text, location, proximity, offset, limit)
	if err != nil {
		handler.logger.Printf("%s unknown error after businessManager.Search %+v", CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	// Convert results to api model
	result := api.GetUserBusinessesSearchResponse{}
	for _, v := range businesses {
		result.Businesses = append(result.Businesses, apiUtils.ConvertBusinessToShortApiModel(&v))
	}

	c.JSON(200, result)
	return
}

// Handles get business info request
// Requires businessId path parameter
func (handler *UserHandlers) getBusiness(c *gin.Context) {
	businessId := c.Param("businessId")

	// Execute the query, handle errors
	businesses, err := handler.businessManager.GetById(businessId, true)
	if err != nil {
		handler.logger.Printf("%s unknown error after businessManager.Search %+v", CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	// Respond with business data
	response := apiUtils.ConvertBusinessToApiModel(businesses, businesses.ItemDefinitions, businesses.MenuImages)
	c.JSON(200, response)
}

// Connects handler to gin router
func (handler *UserHandlers) Connect(rg *gin.RouterGroup) {
	cards := rg.Group("/cards")
	{
		cards.GET("", handler.getUserCards)
		handler.localCardHandlers.Connect(cards.Group("/local"))
		handler.virtualCardHandlers.Connect(cards.Group("/virtual"))
	}
	businesses := rg.Group("/businesses")
	{
		businesses.GET("", handler.getSearchBusinesses)
		businesses.GET("/:businessId", handler.getBusiness)
	}
}

//		UserVirtualCardHandlers

type UserVirtualCardHandlers struct {
	virtualCardManager            VirtualCardManager
	transactionManager            TransactionManager
	itemDefinitionManager         ItemDefinitionManager
	userAuthorizedAcessor         UserAuthorizedAccessor
	authorizedTransactionAccessor AuthorizedTransactionAccessor
	logger                        *log.Logger
}

// Not a request handler. Tries to find virtual card from business with businessId for user.
// Responds with an appropriate error or returns the card.
func (handler *UserVirtualCardHandlers) getVirtualCardOfUser(c *gin.Context,
	user *database.User, businessId string) *database.VirtualCard {

	// Get virtual card, handle errors
	cardTmp, err := handler.userAuthorizedAcessor.Get(user, &database.VirtualCard{
		Business: &database.Business{PublicId: businessId}})

	if err == ErrNotFound {
		c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND, Message: "VIRTUAL_CARD_NOT_FOUND"})
		return nil
	} else if err != nil {
		handler.logger.Printf("%s unknown error after userAuthorizedAcessor.Get: %+v", CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return nil
	}

	return cardTmp.(*database.VirtualCard)
}

// Handles add virtual card request
// Requires businessId path parameter
func (handler *UserVirtualCardHandlers) postCard(c *gin.Context) {
	businessId := c.Param("businessId")

	// Get user from context (should be inserted by authMiddleware)
	user := getUserFromContext(handler.logger, c)
	if user == nil {
		return
	}

	// Create virtual card, handle errors
	_, err := handler.virtualCardManager.Create(user, businessId)
	if err == ErrVirtualCardAlreadyExists {
		c.JSON(409, api.DefaultResponse{Status: api.ALREADY_EXISTS})
		return
	} else if err == ErrNoSuchBusiness {
		c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND, Message: "BUSINESS_NOT_FOUND"})
		return
	} else if err != nil {
		handler.logger.Printf("%s unknown error after virutalCardManager.Create: %+v", CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	c.JSON(201, api.DefaultResponse{Status: api.CREATED})
}

// Handles delete virtual card request
// Requires businessId path parameter
func (handler *UserVirtualCardHandlers) deleteCard(c *gin.Context) {
	businessId := c.Param("businessId")

	// Get user from context (should be inserted by authMiddleware)
	user := getUserFromContext(handler.logger, c)
	if user == nil {
		return
	}

	// Get virtual card of user
	virtualCard := handler.getVirtualCardOfUser(c, user, businessId)
	if virtualCard == nil {
		return
	}

	// Create virtual card, handle errors
	err := handler.virtualCardManager.Remove(virtualCard)
	if err == ErrNoSuchVirtualCard {
		c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND, Message: "VIRTUAL_CARD_NOT_FOUND"})
		return
	} else if err != nil {
		handler.logger.Printf("%s unknown error after virutalCardManager.Remove: %+v", CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	c.JSON(200, api.DefaultResponse{Status: api.OK})
}

// Handles get virtual card info request
// Requires businessId path parameter
func (handler *UserVirtualCardHandlers) getCard(c *gin.Context) {
	businessId := c.Param("businessId")

	// Get user from context (should be inserted by authMiddleware)
	user := getUserFromContext(handler.logger, c)
	if user == nil {
		return
	}

	// Get virtual card, handle errors
	// TODO replace with GetOwnedItems
	cardTmp, err := handler.userAuthorizedAcessor.GetAll(user,
		&database.VirtualCard{Business: &database.Business{PublicId: businessId}},
		[]string{"Business", "Business.ItemDefinitions", "Business.MenuImages",
			"OwnedItems", "OwnedItems.ItemDefinition"})

	if err != nil && err != ErrNotFound {
		handler.logger.Printf("%s unknown error after userAuthorizedAcessor.Get: %+v", CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	} else if err == ErrNotFound {
		c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND, Message: "VIRTUAL_CARD_NOT_FOUND"})
		return
	}

	virtualCard := cardTmp[0].(*database.VirtualCard)

	// Convert data, return response
	var ownedItems []api.OwnedItemApiModel
	for _, v := range virtualCard.OwnedItems {
		ownedItems = append(ownedItems, api.OwnedItemApiModel{
			PublicId:     v.PublicId,
			DefinitionId: v.ItemDefinition.PublicId,
		})
	}

	response := api.GetUserVirtualCardResponse{
		Points: int32(virtualCard.Points),
		BusinessDetails: apiUtils.ConvertBusinessToApiModel(
			virtualCard.Business,
			virtualCard.Business.ItemDefinitions,
			virtualCard.Business.MenuImages),
		OwnedItems: ownedItems,
	}
	c.JSON(200, response)
}

// Handles buy item request
// Requires businessId (matches the virtual card) and itemDefinitionId (public id of item to buy) path parameters
func (handler *UserVirtualCardHandlers) postItemDefinition(c *gin.Context) {
	businessId := c.Param("businessId")
	itemDefinitionId := c.Param("itemDefinitionId")

	// Get user from context (should be inserted by authMiddleware)
	user := getUserFromContext(handler.logger, c)
	if user == nil {
		return
	}

	// Get virtual card of user
	virtualCard := handler.getVirtualCardOfUser(c, user, businessId)
	if virtualCard == nil {
		return
	}

	// Pass to manager, handle errors
	ownedItem, err := handler.virtualCardManager.BuyItem(virtualCard, itemDefinitionId)
	if err == ErrNoSuchItemDefinition {
		c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND})
		return
	} else if err == ErrWithdrawnItem {
		c.JSON(401, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "ITEM_WITHDRAWN"})
		return
	} else if err == ErrUnavailableItem {
		c.JSON(401, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "ITEM_UNAVAILABLE"})
		return
	} else if err == ErrBeforeStartDate {
		c.JSON(401, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "BEFORE_START_DATE"})
		return
	} else if err == ErrAfterEndDate {
		c.JSON(401, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "AFTER_END_DATE"})
		return
	} else if err != nil {
		handler.logger.Printf("%s unknown error after virtualCardManager.BuyItem : %+v", CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	c.JSON(201, api.PostUserVirtualCardItemResponse{
		ItemId: ownedItem.PublicId,
	})
}

// Handles delete (return) item request
// Requires businessId (matches the virtual card) and itemId (public id of ownedItem to return) path parameters
func (handler *UserVirtualCardHandlers) deleteItem(c *gin.Context) {
	businessId := c.Param("businessId")
	itemId := c.Param("itemId")

	// Get user from context (should be inserted by authMiddleware)
	user := getUserFromContext(handler.logger, c)
	if user == nil {
		return
	}

	// Get virtual card of user
	virtualCard := handler.getVirtualCardOfUser(c, user, businessId)
	if virtualCard == nil {
		return
	}

	// Get owned item, handle errors
	ownedItem, err := handler.userAuthorizedAcessor.Get(user, &database.OwnedItem{
		PublicId:      itemId,
		VirtualCardId: virtualCard.ID,
	})

	if err == ErrNoAccess || err == ErrNotFound {
		c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND})
		return
	} else if err != nil {
		handler.logger.Printf("%s unknown error after handler.userAuthorizedAcessor.Get in deleteItem: %+v",
			CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	// Return item, handle errors
	err = handler.virtualCardManager.ReturnItem(ownedItem.(*database.OwnedItem))
	if err == ErrItemCantBeReturned {
		c.JSON(401, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "ITEM_CANNOT_BE_RETURNED"})
		return
	} else if err != nil {
		handler.logger.Printf("%s unknown error after handler.virtualCardManager.ReturnItem in deleteItem: %+v",
			CallerFilename(), err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	c.JSON(200, api.DefaultResponse{Status: api.OK})
}

// Handles start transaction request
// Requires businessId (matches the virtual card) path parameter
func (handler *UserVirtualCardHandlers) postTransaction(c *gin.Context) {
	businessId := c.Param("businessId")

	// Parse request
	req := api.PostUserVirtualCardTransactionRequest{}
	if err := c.BindJSON(&req); err != nil {
		handler.logger.Printf("failed to parse in postTransaction %+v", err)
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST})
		return
	}

	// Get user from context (should be inserted by authMiddleware)
	user := getUserFromContext(handler.logger, c)
	if user == nil {
		return
	}

	// Get virtual card of user
	virtualCard := handler.getVirtualCardOfUser(c, user, businessId)
	if virtualCard == nil {
		return
	}

	// Get items owned by card and present in the request
	ownedItems, err := handler.virtualCardManager.FilterOwnedItems(virtualCard, req.ItemIds)
	if err != nil {
		handler.logger.Printf("unknown error virtualCardManager.FilterOwnedItems in postTransaction %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	// TODO deduplication
	if len(ownedItems) != len(req.ItemIds) {
		c.JSON(401, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "INVALID_ITEM"})
		return
	}

	// Start transaction, handle errors
	transaction, err := handler.transactionManager.Start(virtualCard, ownedItems)
	if err == ErrInvalidItem {
		// TODO how to identify the item?
		c.JSON(401, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "INVALID_ITEM"})
		return
	} else if err != nil {
		handler.logger.Printf("unknown error transactionManager.Start in postTransaction %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	c.JSON(201, api.PostUserVirtualCardTransactionResponse{
		PublicId: transaction.PublicId,
		Code:     transaction.Code,
	})
}

// Handles get transaction info request
// Requires businessId (matches the virtual card) and transactionCode path parameter
func (handler *UserVirtualCardHandlers) getTransaction(c *gin.Context) {
	_ = c.Param("businessId") // TODO this is unused
	transactionCode := c.Param("transactionCode")

	// Get user from context (should be inserted by authMiddleware)
	user := getUserFromContext(handler.logger, c)
	if user == nil {
		return
	}

	// TODO (code,user) and (code,businses) has to be unique, both separately and at once
	// Get transaction data, handle errors
	transaction, err := handler.authorizedTransactionAccessor.GetForUser(user, transactionCode)
	if err == ErrNotFound || err == ErrNoAccess {
		c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND})
		return
	} else if err != nil {
		handler.logger.Printf("unknown error authorizedTransactionAccessor.GetForUser in getTransaction %+v",
			err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	// Convert transaction data to HTTP response
	var itemActions []api.ItemActionApiModel
	for _, v := range transaction.TransactionDetails {
		itemActions = append(itemActions, api.ItemActionApiModel{
			ItemId: v.OwnedItem.PublicId,
			Action: apiUtils.ConvertDbItemAction(v.Action),
		})
	}

	c.JSON(200, api.GetUserVirtualCardTransactionResponse{
		PublicId:    transaction.PublicId,
		State:       apiUtils.ConvertDbTransactionState(transaction.State),
		AddedPoints: int32(transaction.AddedPoints),
		ItemActions: itemActions,
	})
}

func (handler *UserVirtualCardHandlers) Connect(rg *gin.RouterGroup) {
	card := rg.Group("/:businessId")
	{
		card.POST("", handler.postCard)
		card.DELETE("", handler.deleteCard)
		card.GET("", handler.getCard)

		card.POST("/itemsDefinitions/:itemDefinitionId", handler.postItemDefinition)

		card.DELETE("/items/:itemId", handler.deleteItem)

		transactions := card.Group("/transactions")
		{
			transactions.POST("", handler.postTransaction)
			transactions.GET("/:transactionCode", handler.postTransaction)
		}
	}
}

//		UserLocalCardHandlers

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
	user := getUserFromContext(handler.logger, c)
	if user == nil {
		return
	}

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
			ImageUrl: v.ImageUrl,
		})
	}

	c.JSON(200, api.GetUserLocalCardTypesResponse{
		Types: types,
	})
}

// Handles delete card request
// Requires publicId URL path parameter
func (handler *UserLocalCardHandlers) deleteCard(c *gin.Context) {
	cardId := c.Param("cardId")

	// Get user from context (should be inserted by authMiddleware)
	user := getUserFromContext(handler.logger, c)
	if user == nil {
		return
	}

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
