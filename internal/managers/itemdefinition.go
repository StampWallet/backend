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
			return fmt.Errorf("fileStorageService.CreateStub returned an error: %w", err)
		}

		itemDefinition = ItemDefinition{
			PublicId:    shortuuid.New(),
			BusinessId:  business.ID,
			Name:        details.Name,
			Price:       *details.Price,
			Description: details.Description,
			ImageId:     imageFile.PublicId, // TODO: expect PublicId here?
			StartDate:   sql.NullTime{Time: *details.StartDate, Valid: true},
			EndDate:     sql.NullTime{Time: *details.EndDate, Valid: true},
			MaxAmount:   *details.MaxAmount,
			Available:   *details.Available,
			Withdrawn:   false,
		}

		result := db.Create(&itemDefinition)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Create returned an error: %w", err)
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
		item.StartDate = sql.NullTime{Time: *details.StartDate, Valid: true}
	}
	if details.EndDate != nil {
		item.EndDate = sql.NullTime{Time: *details.EndDate, Valid: true}
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
			return fmt.Errorf("db.Find(ItemDefinition) returned an error %w", err)
		}

		item.Withdrawn = true
		item.Available = false

		execDb := db.Exec(`UPDATE virtual_cards AS vc
				SET points = vc.points + t.points
			FROM		
				(SELECT oi.virtual_card_id as vid, sum(itd.price) as points
				FROM owned_items oi 
					JOIN item_definitions AS itd ON itd.id = oi.definition_id 
				WHERE oi.definition_id=? AND oi.used is NULL AND oi.status='OWNED'
				GROUP BY oi.virtual_card_id) AS t
			WHERE
				vc.id = t.vid`, item.ID)

		if err := execDb.GetError(); err != nil {
			return fmt.Errorf("failed to update virtual cards in WithdrawItem: %w", err)
		}

		execDb = db.Exec(`UPDATE owned_items SET status=? WHERE definition_id=? AND status=?`,
			OwnedItemStatusWithdrawn, item.ID, OwnedItemStatusOwned)

		if err := execDb.GetError(); err != nil {
			return fmt.Errorf("failed to disable owned items in WithdrawItem: %w", err)
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
