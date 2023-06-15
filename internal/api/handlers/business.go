package api

import (
	"log"

	"github.com/gin-gonic/gin"

	api "github.com/StampWallet/backend/internal/api/models"
	"github.com/StampWallet/backend/internal/api/utils"
	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/database/accessors"
	. "github.com/StampWallet/backend/internal/managers"
	"github.com/StampWallet/backend/internal/services"
)

type BusinessHandlers struct {
	businessManager       BusinessManager
	transactionManager    TransactionManager
	itemDefinitionManager ItemDefinitionManager

	userAuthorizedAcessor         UserAuthorizedAccessor
	businessAuthorizedAccessor    BusinessAuthorizedAccessor
	authorizedTransactionAccessor AuthorizedTransactionAccessor

	itemDefinitionHandlers *ItemDefinitionHandlers

	logger *log.Logger
}

func CreateBusinessHandlers(
	businessManager BusinessManager, transactionManager TransactionManager,
	itemDefinitionManager ItemDefinitionManager,
	userAuthorizedAcessor UserAuthorizedAccessor, businessAuthorizedAccessor BusinessAuthorizedAccessor,
	authorizedTransactionAccessor AuthorizedTransactionAccessor,
	logger *log.Logger) *BusinessHandlers {

	return &BusinessHandlers{
		businessManager:       businessManager,
		transactionManager:    transactionManager,
		itemDefinitionManager: itemDefinitionManager,

		userAuthorizedAcessor:         userAuthorizedAcessor,
		businessAuthorizedAccessor:    businessAuthorizedAccessor,
		authorizedTransactionAccessor: authorizedTransactionAccessor,

		itemDefinitionHandlers: &ItemDefinitionHandlers{
			itemDefinitionManager:      itemDefinitionManager,
			userAuthorizedAcessor:      userAuthorizedAcessor,
			businessAuthorizedAccessor: businessAuthorizedAccessor,
			logger:                     services.NewPrefix(logger, "ItemDefinitionHandlers"),
		},

		logger: logger,
	}
}

// Gets business owned by user
func (handler *BusinessHandlers) getBusinessOfUser(user *User, c *gin.Context) *Business {
	// Get user's business
	businessTmp, err := handler.userAuthorizedAcessor.Get(user, &Business{})

	if err != nil && err != ErrNotFound {
		handler.logger.Printf("failed to handler.userAuthorizedAcessor.Get in getAccountInfo %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return nil
	} else if businessTmp == nil || err == ErrNotFound {
		c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND})
		return nil
	}

	return businessTmp.(*Business)
}

// Gets user (usually inserted into the context by AuthMiddleware)
// and business owned by user from request context.
func (handler *BusinessHandlers) getUserAndBusiness(c *gin.Context) (*User, *Business) {
	// Get user from context
	user := getUserFromContext(handler.logger, c)
	if user == nil {
		return nil, nil
	}

	// Get user's business
	business := handler.getBusinessOfUser(user, c)
	if business == nil {
		return nil, nil
	}

	return user, business
}

// Handles business account creation request
func (handler *BusinessHandlers) postAccount(c *gin.Context) {
	// Parse request body
	req := api.PostBusinessAccountRequest{}
	if err := c.BindJSON(&req); err != nil {
		handler.logger.Printf("failed to parse in postAccount %+v", err)
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST})
		return
	}

	// Get user from context (should be inserted by authMiddleware)
	user := getUserFromContext(handler.logger, c)
	if user == nil {
		return
	}

	// Parse coordinates
	coordinates, err := GPSCoordinatesFromString(req.GpsCoordinates)
	if err != nil {
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "INVALID_GPS_COORDINATES"})
		return
	}

	// Pass request to businessManager
	//TODO add data validation
	business, err := handler.businessManager.Create(user, &BusinessDetails{
		Name:           req.Name,
		Description:    req.Description,
		Address:        req.Address,
		NIP:            req.Nip,
		KRS:            req.Krs,
		REGON:          req.Regon,
		OwnerName:      req.OwnerName,
		GPSCoordinates: coordinates,
	})

	// Handle errors, send response
	if err != nil {
		handler.logger.Printf("failed to businessManager.Create in postAccount %+v", err)
		if err == ErrBusinessAlreadyExists {
			c.JSON(409, api.DefaultResponse{Status: api.ALREADY_EXISTS})
			return
		} else {
			c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
			return
		}
	}

	c.JSON(201, api.PostBusinessAccountResponse{
		PublicId:      business.PublicId,
		BannerImageId: business.BannerImageId,
		IconImageId:   business.IconImageId,
	})
}

// Handles business info retrieval
func (handler *BusinessHandlers) getAccountInfo(c *gin.Context) {
	// Get user and business of user. getUserAndBusiness sends HTTP errors, so we can just quit
	// if business or user is not available
	user, business := handler.getUserAndBusiness(c)
	if user == nil || business == nil {
		return
	}

	// Get MenuItems of business
	menuItems, err := handler.businessAuthorizedAccessor.GetAll(business, &MenuImage{})
	if err != nil {
		handler.logger.Printf("failed to handler.businessAuthorizedAccessor.Get MenuItem in getAccountInfo %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	// Get ItemDefinitions of business
	itemDefinitionsTmp, err := handler.businessAuthorizedAccessor.GetAll(business, &ItemDefinition{Withdrawn: false})
	if err != nil {
		handler.logger.Printf("failed to handler.businessAuthorizedAccessor.Get ItemDefinition in getAccountInfo %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	// Convert MenuItems
	var menuImageIds []string
	for _, v := range menuItems {
		menuImageIds = append(menuImageIds, v.(*MenuImage).FileId)
	}

	// Convert ItemDefinitions
	var itemDefinitions []api.ItemDefinitionApiModel
	for _, v := range itemDefinitionsTmp {
		itemDefinitions = append(itemDefinitions, apiUtils.ConvertItemDefinitionToApiModel(v.(*ItemDefinition)))
	}

	c.JSON(200, api.GetBusinessAccountResponse{
		PublicId:        business.PublicId,
		Name:            business.Name,
		Address:         business.Address,
		GpsCoordinates:  business.GPSCoordinates.ToString(),
		BannerImageId:   business.BannerImageId,
		IconImageId:     business.IconImageId,
		MenuImageIds:    menuImageIds,
		ItemDefinitions: itemDefinitions,
		Nip:             business.NIP,
		Krs:             business.KRS,
		Regon:           business.REGON,
		OwnerName:       business.OwnerName,
		Description:     business.Description,
	})
}

// Handles business info change request
func (handler *BusinessHandlers) patchAccountInfo(c *gin.Context) {
	// Parse request body
	req := api.PatchBusinessAccountRequest{}
	if err := c.BindJSON(&req); err != nil {
		handler.logger.Printf("failed to parse in postAccount %+v", err)
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST})
		return
	}

	var nameToChange *string
	var descriptionToChange *string

	if req.Name != "" {
		nameToChange = &req.Name
	}

	if req.Description != "" {
		descriptionToChange = &req.Description
	}

	// Make sure that the request is correct - either Name or Description has to be not empty
	if nameToChange == nil && descriptionToChange == nil {
		c.JSON(401, api.DefaultResponse{Status: api.INVALID_REQUEST})
		return
	}

	// Get user and business of user. getUserAndBusiness sends HTTP errors, so we can just quit
	// if business or user is not available
	user, business := handler.getUserAndBusiness(c)
	if user == nil || business == nil {
		return
	}

	// Send to manager, handle errors, send response
	_, err := handler.businessManager.ChangeDetails(business, &ChangeableBusinessDetails{
		Name:        nameToChange,
		Description: descriptionToChange,
	})

	if err != nil {
		handler.logger.Printf("failed to handler.businessManager.ChangeDetails in patchAccountInfo %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	c.JSON(200, api.DefaultResponse{Status: api.OK})
}

// Handles transaction get request from business account.
// Requires {transactionCode} URL path parameter
func (handler *BusinessHandlers) getTransaction(c *gin.Context) {
	transactionCode := c.Param("transactionCode")

	// Get user and business of user. getUserAndBusiness sends HTTP errors, so we can just quit
	// if business or user is not available
	user, business := handler.getUserAndBusiness(c)
	if user == nil || business == nil {
		return
	}

	// Get transaction
	transaction, err := handler.authorizedTransactionAccessor.GetForBusiness(business, transactionCode)
	if err != nil {
		handler.logger.Printf("failed to handler.authorizedTransactionAccessor.GetForBusiness in getTransaction %+v", err)
		if err == ErrNoAccess {
			c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND})
			return
		} else if err != nil {
			c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
			return
		}
	}

	//NOTE client will have to cache item definitions
	// Currently, I don't think the client has to know anything about OwnedItem
	// besides its ItemDefinition and PublicId.
	// Used and Status are used only on the backend to enforce some rules.
	// TODO add "LastUpdated" fields or headers to API Models for tansactons/itemDefinitions/businesses
	var transactionItems []api.TransactionItemDetailApiModel
	for _, v := range transaction.TransactionDetails {
		transactionItems = append(transactionItems, api.TransactionItemDetailApiModel{
			PublicId:         v.OwnedItem.PublicId,
			ItemDefinitionId: v.OwnedItem.ItemDefinition.PublicId,
		})
	}

	c.JSON(200, api.GetBusinessTransactionResponse{
		PublicId:      transaction.PublicId,
		VirtualCardId: int32(transaction.VirtualCardId),
		State:         apiUtils.ConvertDbTransactionState(transaction.State),
		Items:         transactionItems,
	})
}

// Handles transaction finalize request
// Requires {transactionCode} URL path parameter
func (handler *BusinessHandlers) postTransaction(c *gin.Context) {
	transactionCode := c.Param("transactionCode")

	// Parse request body
	req := api.PostBusinessTransactionRequest{}
	if err := c.BindJSON(&req); err != nil {
		handler.logger.Printf("failed to parse in postTransaction %+v", err)
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST})
		return
	}

	// Get user and business of user. getUserAndBusiness sends HTTP errors, so we can just quit
	// if business or user is not available
	user, business := handler.getUserAndBusiness(c)
	if user == nil || business == nil {
		return
	}

	// Get transaction
	transaction, err := handler.authorizedTransactionAccessor.GetForBusiness(business, transactionCode)
	if err != nil {
		handler.logger.Printf("failed to handler.authorizedTransactionAccessor.GetForBusiness in postTransaction %+v", err)
		if err == ErrNoAccess {
			c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND})
			return
		} else if err != nil {
			c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
			return
		}
	}

	// Convert request
	var itemMap map[string]*OwnedItem = map[string]*OwnedItem{}
	for _, v := range transaction.TransactionDetails {
		itemMap[v.OwnedItem.PublicId] = v.OwnedItem
	}

	var itemActions []ItemWithAction
	for _, v := range req.ItemActions {
		item, ok := itemMap[v.ItemId]
		if !ok {
			handler.logger.Printf("unknown item in postTransaction %s", v.ItemId)
			c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST, Message: "UNKNOWN_ITEM"})
			return
		}
		itemActions = append(itemActions, ItemWithAction{
			Item:   item,
			Action: apiUtils.ConvertApiItemAction(v.Action),
		})
	}

	// Send data to manager, handle errors
	_, err = handler.transactionManager.Finalize(transaction, itemActions, uint64(req.AddedPoints))
	if err != nil {
		handler.logger.Printf("failed to handler.transactionManager.Finalize in postTransaction %+v", err)
		if err == ErrInvalidItem {
			c.JSON(400, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
			return
		} else {
			c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
			return
		}
	}

	c.JSON(200, api.DefaultResponse{Status: api.OK})
}

// Handles menu image add request
func (handler *BusinessHandlers) postMenuImage(c *gin.Context) {
	// Get user and business of user. getUserAndBusiness sends HTTP errors, so we can just quit
	// if business or user is not available
	user, business := handler.getUserAndBusiness(c)
	if user == nil || business == nil {
		return
	}

	menuImage, err := handler.businessManager.AddMenuImage(user, business)
	if err == ErrTooManyMenuImages {
		c.JSON(400, api.DefaultResponse{
			Status:  api.INVALID_REQUEST,
			Message: "TOO_MANY_IMAGES",
		})
		return
	} else if err != nil {
		handler.logger.Printf("failed to handler.businessManager.AddMenuImage in postMenuImage %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	c.JSON(200, api.PostBusinessAccountMenuImageResponse{
		ImageId: menuImage.FileId,
	})
}

// Handles menu image remove request
// Requires {menuImageId} URL path parameter
func (handler *BusinessHandlers) deleteMenuImage(c *gin.Context) {
	menuImageId := c.Param("menuImageId")

	// Get user and business of user. getUserAndBusiness sends HTTP errors, so we can just quit
	// if business or user is not available
	user, business := handler.getUserAndBusiness(c)
	if user == nil || business == nil {
		return
	}

	menuImageTmp, err := handler.businessAuthorizedAccessor.Get(business, &MenuImage{FileId: menuImageId})
	if err == ErrNoAccess {
		c.JSON(403, api.DefaultResponse{Status: api.FORBIDDEN})
		return
	} else if err == ErrNotFound {
		c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND})
		return
	} else if err != nil {
		handler.logger.Printf("failed to handler.businessAuthorizedAccessor.Get in deleteMenuImage %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	menuImage := menuImageTmp.(*MenuImage)

	err = handler.businessManager.RemoveMenuImage(menuImage)
	if err != nil {
		handler.logger.Printf("failed to handler.businessManager.RemoveMenuImage in deleteMenuImage %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	c.JSON(200, api.DefaultResponse{Status: api.OK})
}

func (handler *BusinessHandlers) Connect(rg *gin.RouterGroup) {
	rg.POST("/account", handler.postAccount)
	rg.GET("/info", handler.getAccountInfo)
	rg.PATCH("/info", handler.patchAccountInfo)

	menuImages := rg.Group("/menuImages")
	{
		menuImages.POST("/", handler.postMenuImage)
		menuImages.DELETE("/:menuImageId", handler.deleteMenuImage)
	}

	transactions := rg.Group("/transactions")
	{
		transactions.GET("/:transactionCode", handler.getTransaction)
		transactions.POST("/:transactionCode", handler.postTransaction)
	}

	handler.itemDefinitionHandlers.Connect(rg.Group("/itemDefinitions"))
}

// ItemDefinitionHandlers

type ItemDefinitionHandlers struct {
	itemDefinitionManager      ItemDefinitionManager
	userAuthorizedAcessor      UserAuthorizedAccessor
	businessAuthorizedAccessor BusinessAuthorizedAccessor
	logger                     *log.Logger
}

// Gets business owned by user
// TODO refactor - duplicated from BusinessHandlers
func (handler *ItemDefinitionHandlers) getBusinessOfUser(user *User, c *gin.Context) *Business {
	// Get user's business
	businessTmp, err := handler.userAuthorizedAcessor.Get(user, &Business{})

	if err != nil && err != ErrNotFound {
		handler.logger.Printf("failed to handler.userAuthorizedAcessor.Get in getAccountInfo %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return nil
	} else if businessTmp == nil || err == ErrNotFound {
		c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND})
		return nil
	}

	return businessTmp.(*Business)
}

// Gets user (usually inserted into the context by AuthMiddleware)
// and business owned by user from request context.
// TODO refactor - duplicated from BusinessHandlers
func (handler *ItemDefinitionHandlers) getUserAndBusiness(c *gin.Context) (*User, *Business) {
	// Get user from context
	user := getUserFromContext(handler.logger, c)
	if user == nil {
		return nil, nil
	}

	// Get user's business
	business := handler.getBusinessOfUser(user, c)
	if business == nil {
		return nil, nil
	}

	return user, business
}

func (handler *ItemDefinitionHandlers) getItemDefinition(c *gin.Context) {
	// Get user and business of user. getUserAndBusiness sends HTTP errors, so we can just quit
	// if business or user is not available
	user, business := handler.getUserAndBusiness(c)
	if user == nil || business == nil {
		return
	}

	itemDefinitionsTmp, err := handler.businessAuthorizedAccessor.GetAll(business, &ItemDefinition{BusinessId: business.ID, Withdrawn: false})
	if err != nil {
		handler.logger.Printf("failed to handler.businessAuthorizedAccessor.GetAll in getItemDefinition%+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	var itemDefinitions []api.ItemDefinitionApiModel
	for _, v := range itemDefinitionsTmp {
		itemDefinitions = append(itemDefinitions, apiUtils.ConvertItemDefinitionToApiModel(v.(*ItemDefinition)))
	}

	c.JSON(200, api.GetBusinessItemDefinitionsResponse{
		ItemDefinitions: itemDefinitions,
	})
}

func (handler *ItemDefinitionHandlers) postItemDefinition(c *gin.Context) {
	user, business := handler.getUserAndBusiness(c)
	if user == nil || business == nil {
		return
	}

	// Parse request body
	req := api.PostBusinessItemDefinitionRequest{}
	if err := c.BindJSON(&req); err != nil {
		handler.logger.Printf("failed to parse in postItemDefinition %+v", err)
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST})
		return
	}

	var price *uint
	var maxAmount *uint

	if req.Price != nil {
		tmp := uint(*req.Price)
		price = &tmp
	}

	if req.MaxAmount != nil {
		tmp := uint(*req.MaxAmount)
		maxAmount = &tmp
	}

	itemDefinition, err := handler.itemDefinitionManager.AddItem(user, business, &ItemDetails{
		Name:        req.Name,
		Price:       price,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		MaxAmount:   maxAmount,
		Available:   &req.Available,
	})

	if err == ErrInvalidItemDetails {
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST})
		return
	} else if err != nil {
		handler.logger.Printf("unknown error returned from itemDefinitionManager.AddItem: %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	c.JSON(201, api.PostBusinessItemDefinitionResponse{PublicId: itemDefinition.PublicId})
}

// Requires {definitionId} URL path parameter
func (handler *ItemDefinitionHandlers) putItemDefinition(c *gin.Context) {
	definitionId := c.Param("definitionId")

	user, business := handler.getUserAndBusiness(c)
	if user == nil || business == nil {
		return
	}

	// Parse request body
	req := api.PutBusinessItemDefinitionRequest{}
	if err := c.BindJSON(&req); err != nil {
		handler.logger.Printf("failed to parse in putItemDefinition %+v", err)
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST})
		return
	}

	itemDefinitionTmp, err := handler.businessAuthorizedAccessor.Get(business, &ItemDefinition{PublicId: definitionId})
	if err == ErrNoAccess || err == ErrNotFound {
		c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND})
		return
	} else if err != nil {
		handler.logger.Printf("failed to handler.businessAuthorizedAccessor.Get in putItemDefinition: %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	var price *uint
	var maxAmount *uint

	if req.Price != nil {
		tmp := uint(*req.Price)
		price = &tmp
	}

	if req.MaxAmount != nil {
		tmp := uint(*req.MaxAmount)
		maxAmount = &tmp
	}

	_, err = handler.itemDefinitionManager.ChangeItemDetails(itemDefinitionTmp.(*ItemDefinition),
		&ItemDetails{
			Name:        req.Name,
			Price:       price,
			Description: req.Description,
			StartDate:   req.StartDate,
			EndDate:     req.EndDate,
			MaxAmount:   maxAmount,
			Available:   &req.Available,
		})

	if err == ErrInvalidItemDetails {
		c.JSON(400, api.DefaultResponse{Status: api.INVALID_REQUEST})
		return
	} else if err != nil {
		handler.logger.Printf("unknown error returned from itemDefinitionManager.AddItem: %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	c.JSON(200, api.DefaultResponse{Status: api.OK})
}

// Requires {definitionId} URL path parameter
func (handler *ItemDefinitionHandlers) deleteItemDefinition(c *gin.Context) {
	definitionId := c.Param("definitionId")

	user, business := handler.getUserAndBusiness(c)
	if user == nil || business == nil {
		return
	}

	itemDefinitionTmp, err := handler.businessAuthorizedAccessor.Get(business, &ItemDefinition{PublicId: definitionId})
	if err == ErrNoAccess || err == ErrNotFound {
		c.JSON(404, api.DefaultResponse{Status: api.NOT_FOUND})
		return
	} else if err != nil {
		handler.logger.Printf("failed to handler.businessAuthorizedAccessor.Get in deleteItemDefinition: %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	_, err = handler.itemDefinitionManager.WithdrawItem(itemDefinitionTmp.(*ItemDefinition))
	if err != nil {
		handler.logger.Printf("failed to itemDefinitionManager.WithdrawItem in deleteItemDefinition: %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	c.JSON(200, api.DefaultResponse{Status: api.OK})
}

func (handler *ItemDefinitionHandlers) Connect(rg *gin.RouterGroup) {
	rg.GET("", handler.getItemDefinition)
	rg.POST("", handler.postItemDefinition)
	rg.PUT("/:definitionId", handler.putItemDefinition)
	rg.DELETE("/:definitionId", handler.deleteItemDefinition)
}
