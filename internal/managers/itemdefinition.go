package managers

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services"
	"github.com/lithammer/shortuuid/v4"
)

type ItemDefinitionManager interface {
	AddItem(user *User, business *Business, details *ItemDetails) (*ItemDefinition, error)
	ChangeItemDetails(item *ItemDefinition, details *ItemDetails) (*ItemDefinition, error)
	WithdrawItem(item *ItemDefinition) (*ItemDefinition, error)
	GetForBusiness(business *Business) ([]ItemDefinition, error)
}

type ItemDetails struct {
	Name        string
	Price       *uint
	Description string
	StartDate   *time.Time
	EndDate     *time.Time
	MaxAmount   *uint
	Available   *bool
}

var ErrInvalidItemDetails = errors.New("Invalid item details received")
var ErrUnknownItem = errors.New("Item id is not recognized")
var ErrInvalidArgs = errors.New("Arguments passed to manager are invalid")
var ErrItemAlreadyWithdrawn = errors.New("Item already withdrawn")

type ItemDefinitionManagerImpl struct {
	baseServices       BaseServices
	fileStorageService FileStorageService
}

func CreateItemDefinitionManagerImpl(baseServices BaseServices, fileStorageService FileStorageService) *ItemDefinitionManagerImpl {
	return &ItemDefinitionManagerImpl{
		baseServices:       baseServices,
		fileStorageService: fileStorageService,
	}
}

// TODO add logs
func (manager *ItemDefinitionManagerImpl) AddItem(user *User, business *Business, details *ItemDetails) (*ItemDefinition, error) {
	var itemDefinition ItemDefinition

	if user.ID != business.OwnerId {
		return nil, ErrInvalidArgs
	}

	err := manager.baseServices.Database.Transaction(func(db GormDB) error {
		imageFile, err := manager.fileStorageService.CreateStub(user)
		if err != nil {
			return fmt.Errorf("fileStorageService.CreateStub returned an error: %+v", err)
		}

		itemDefinition = ItemDefinition{
			PublicId:    shortuuid.New(),
			BusinessId:  business.ID,
			Name:        details.Name,
			Price:       *details.Price,
			Description: details.Description,
			ImageId:     imageFile.PublicId, // TODO: expect PublicId here?
			StartDate:   sql.NullTime{*details.StartDate, true},
			EndDate:     sql.NullTime{*details.EndDate, true},
			MaxAmount:   *details.MaxAmount,
			Available:   *details.Available,
			Withdrawn:   false,
		}

		result := db.Create(&itemDefinition)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Create returned an error: %+v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &itemDefinition, nil
}

func (manager *ItemDefinitionManagerImpl) ChangeItemDetails(item *ItemDefinition, details *ItemDetails) (*ItemDefinition, error) {
	if details.Name != "" {
		item.Name = details.Name
	}
	if details.Price != nil {
		item.Price = *details.Price
	}
	if details.Description != "" {
		item.Description = details.Description
	}
	if details.StartDate != nil {
		item.StartDate = sql.NullTime{*details.StartDate, true}
	}
	if details.EndDate != nil {
		item.EndDate = sql.NullTime{*details.EndDate, true}
	}
	if details.MaxAmount != nil {
		item.MaxAmount = *details.MaxAmount
	}
	if details.Available != nil {
		item.Available = *details.Available
	}

	tx := manager.baseServices.Database.Save(item)
	if err := tx.GetError(); err != nil {
		return nil, err
	}

	return item, nil
}

func (manager *ItemDefinitionManagerImpl) WithdrawItem(item *ItemDefinition) (*ItemDefinition, error) {
	err := manager.baseServices.Database.Transaction(func(db GormDB) error {
		if item.Withdrawn {
			return ErrItemAlreadyWithdrawn
		}

		result := db.
			Preload("OwnedItems").
			Preload("OwnedItems.VirtualCard").
			Find(item, "id = ?", item.ID)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Find(ItemDefinition) returned an error %+v", err)
		}

		item.Withdrawn = true
		item.Available = false
		// TODO: not scalable
		for _, ownedItem := range item.OwnedItems {
			if ownedItem.Status != OwnedItemStatusOwned || ownedItem.Used.Valid {
				continue
			}
			ownedItem.Status = OwnedItemStatusWithdrawn
			ownedItem.VirtualCard.Points += item.Price

			result = db.Save(&ownedItem)
			if err := result.GetError(); err != nil {
				return fmt.Errorf("db.Save(ownedItem) returned an error %+v", err)
			}

			result = db.Save(&ownedItem.VirtualCard)
			if err := result.GetError(); err != nil {
				return fmt.Errorf("db.Save(ownedItem.VirtualCard) returned an error %+v", err)
			}
		}
		db.Save(item)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (manager *ItemDefinitionManagerImpl) GetForBusiness(business *Business) ([]ItemDefinition, error) {
	result := manager.baseServices.Database.
		Preload("ItemDefinitions").
		Find(&business, "id = ?", business.ID)
	if err := result.GetError(); err != nil {
		return nil, err
	}
	return business.ItemDefinitions, nil
}
