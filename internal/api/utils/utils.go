package apiUtils

import (
	"fmt"
	"time"

	"github.com/StampWallet/backend/internal/api/models"
	"github.com/StampWallet/backend/internal/database"
)

// on one hand, this is not necessary. both enums have the same values
// on the other, this could change, and I would prefer a crash if that ever happens
// it seems enums in go provide exactly zero type safety - just another type alias for string
func ConvertApiTransactionState(arg api.TransactionStateEnum) database.TransactionStateEnum {
	if arg == api.STARTED {
		return database.TransactionStateStarted
	} else if arg == api.PROCESSING {
		return database.TransactionStateProcesing
	} else if arg == api.FINISHED {
		return database.TransactionStateFinished
	} else if arg == api.EXPIRED {
		return database.TransactionStateExpired
	} else {
		panic(fmt.Errorf("unkown api.TransactionStateEnum enum valule - cannot map to database.TransactionStateEnum %+v", arg))
	}
}

func ConvertDbTransactionState(arg database.TransactionStateEnum) api.TransactionStateEnum {
	if arg == database.TransactionStateStarted {
		return api.STARTED
	} else if arg == database.TransactionStateProcesing {
		return api.PROCESSING
	} else if arg == database.TransactionStateFinished {
		return api.FINISHED
	} else if arg == database.TransactionStateExpired {
		return api.EXPIRED
	} else {
		panic(fmt.Errorf("unkown api.TransactionStateEnum enum valule - cannot map to database.TransactionStateEnum %+v", arg))
	}
}

func ConvertApiItemAction(arg api.ItemActionTypeEnum) database.ActionTypeEnum {
	if arg == api.REDEEMED {
		return database.RedeemedActionType
	} else if arg == api.NO_ACTION {
		return database.NoActionType
	} else if arg == api.RECALLED {
		return database.RecalledActionType
	} else if arg == api.CANCELLED {
		return database.CancelledActionType
	} else {
		panic(fmt.Errorf("unkown api.TransactionStateEnum enum valule - cannot map to database.TransactionStateEnum %+v", arg))
	}
}

func ConvertDbItemAction(arg database.ActionTypeEnum) api.ItemActionTypeEnum {
	if arg == database.NoActionType {
		return api.NO_ACTION
	} else if arg == database.RedeemedActionType {
		return api.REDEEMED
	} else if arg == database.RecalledActionType {
		return api.RECALLED
	} else if arg == database.CancelledActionType {
		return api.CANCELLED
	} else {
		panic(fmt.Errorf("unkown api.TransactionStateEnum enum valule - cannot map to database.TransactionStateEnum %+v", arg))
	}
}

// Converts ItemDefinition from database model to api model
func ConvertItemDefinitionToApiModel(itd *database.ItemDefinition) api.ItemDefinitionApiModel {
	var sd *time.Time
	if itd.StartDate.Valid {
		sd = &itd.StartDate.Time
	}

	var ed *time.Time
	if itd.EndDate.Valid {
		ed = &itd.EndDate.Time
	}

	price := int32(itd.Price)
	maxAmount := int32(itd.MaxAmount)

	return api.ItemDefinitionApiModel{
		PublicId:    itd.PublicId,
		Name:        itd.Name,
		Price:       &price,
		Description: itd.Description,
		ImageId:     itd.ImageId,
		StartDate:   sd,
		EndDate:     ed,
		MaxAmount:   &maxAmount,
		Available:   itd.Available,
	}
}

// Converts database.Business to api.ShortBusinessDetailsApiModel
// Most data is lost in conversion - api.ShortBusinessDetailsApiModel does not contain all
// data from model
func ConvertBusinessToShortApiModel(business *database.Business) api.ShortBusinessDetailsApiModel {
	return api.ShortBusinessDetailsApiModel{
		PublicId:       business.PublicId,
		Name:           business.Name,
		Description:    business.Description,
		GpsCoordinates: business.GPSCoordinates.ToString(),
		BannerImageId:  business.BannerImageId,
		IconImageId:    business.IconImageId,
	}
}

// Converts database.Business to api.PublicBusinessDetailsAPIModel
// Some data is lost in conversion - api.PublicBusinessDetailsAPIModel does not contain all
// data from model
func ConvertBusinessToApiModel(business *database.Business,
	itemDefinitions []database.ItemDefinition,
	menuImages []database.MenuImage) api.PublicBusinessDetailsApiModel {

	var menuImageIds []string
	for _, v := range menuImages {
		menuImageIds = append(menuImageIds, v.FileId)
	}

	var itemDefinitionsApi []api.ItemDefinitionApiModel
	for _, v := range itemDefinitions {
		itemDefinitionsApi = append(itemDefinitionsApi, ConvertItemDefinitionToApiModel(&v))
	}

	return api.PublicBusinessDetailsApiModel{
		PublicId:        business.PublicId,
		Name:            business.Name,
		GpsCoordinates:  business.GPSCoordinates.ToString(),
		BannerImageId:   business.BannerImageId,
		Description:     business.Description,
		IconImageId:     business.IconImageId,
		MenuImageIds:    menuImageIds,
		Address:         business.Address,
		ItemDefinitions: itemDefinitionsApi,
	}
}
