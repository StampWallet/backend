package api

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"

	api "github.com/StampWallet/backend/internal/api/models"
	"github.com/StampWallet/backend/internal/api/utils"
	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/database/accessors"
	. "github.com/StampWallet/backend/internal/managers"
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

// Converts ItemDefinition from database model to api model
func convertItemDefinitionToApiModel(itd *ItemDefinition) api.ItemDefinitionApiModel {
	var sd *time.Time
	if itd.StartDate.Valid {
		sd = &itd.StartDate.Time
	}

	var ed *time.Time
	if itd.EndDate.Valid {
		ed = &itd.EndDate.Time
	}

	return api.ItemDefinitionApiModel{
		PublicId:    itd.PublicId,
		Name:        itd.Name,
		Price:       int32(itd.Price),
		Description: itd.Description,
		ImageId:     itd.ImageId,
		StartDate:   sd,
		EndDate:     ed,
		MaxAmount:   int32(itd.MaxAmount),
		Available:   itd.Available,
	}
}

// Gets business owned by user
func (handler *BusinessHandlers) getBusinessOfUser(user *User, c *gin.Context) *Business {
	// Get user's business
	businessTmp, err := handler.userAuthorizedAcessor.Get(user, &Business{})

	if err != nil {
		handler.logger.Printf("failed to handler.userAuthorizedAcessor.Get in getAccountInfo %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return nil
	} else if businessTmp == nil {
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
		if err == BusinessAlreadyExists {
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
	menuItems, err := handler.businessAuthorizedAccessor.GetAll(business, &MenuItem{})
	if err != nil {
		handler.logger.Printf("failed to handler.businessAuthorizedAccessor.Get MenuItem in getAccountInfo %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	// Get ItemDefinitions of business
	itemDefinitionsTmp, err := handler.businessAuthorizedAccessor.GetAll(business, &ItemDefinition{})
	if err != nil {
		handler.logger.Printf("failed to handler.businessAuthorizedAccessor.Get ItemDefinition in getAccountInfo %+v", err)
		c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
		return
	}

	// Convert MenuItems
	var menuImageIds []string
	for _, v := range menuItems {
		menuImageIds = append(menuImageIds, v.(*MenuItem).FileId)
	}

	// Convert ItemDefinitions
	var itemDefinitions []api.ItemDefinitionApiModel
	for _, v := range itemDefinitionsTmp {
		itemDefinitions = append(itemDefinitions, convertItemDefinitionToApiModel(v.(*ItemDefinition)))
	}

	c.JSON(201, api.GetBusinessAccountResponse{
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
	var itemMap map[string]*OwnedItem
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
		if err == InvalidItem {
			c.JSON(400, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
			return
		} else {
			c.JSON(500, api.DefaultResponse{Status: api.UNKNOWN_ERROR})
			return
		}
	}

	c.JSON(200, api.DefaultResponse{Status: api.OK})
}

func (handler *BusinessHandlers) postMenuImage(c *gin.Context) {

}

func (handler *BusinessHandlers) deleteMenuImage(c *gin.Context) {

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
	rg.GET("", handler.getItemDefinition)
	rg.POST("", handler.postItemDefinition)
	rg.PATCH("/:definitionId", handler.patchItemDefinition)
	rg.DELETE("/:definitionId", handler.deleteItemDefinition)
}
