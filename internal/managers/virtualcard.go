package managers

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services"
	"github.com/lithammer/shortuuid/v4"
	"gorm.io/gorm"
)

var (
	ErrVirtualCardAlreadyExists = errors.New("Virtual card already exists")
	ErrNoSuchVirtualCard        = errors.New("Virtual card not found")
	ErrNoSuchItemDefinition     = errors.New("Item definition not found")
	ErrAboveMaxAmount           = errors.New("Attempt to buy item above max item amount")
	ErrNotEnoughPoints          = errors.New("Attempt to buy item with not enough points")
	ErrBeforeStartDate          = errors.New("Attempt to buy item before start date")
	ErrAfterEndDate             = errors.New("Attempt to buy item after end date")
	ErrUnavailableItem          = errors.New("Attempt to buy an unavailable item")
	ErrWithdrawnItem            = errors.New("Attempt to buy a withdrawn item")
	ErrItemCantBeReturned       = errors.New("Item can't be returned")
)

type VirtualCardManager interface {
	// Creates virtual card of business for user
	// businessId passed as string - caller is not required to have "access" to a business object
	Create(user *User, businessId string) (*VirtualCard, error)

	// Removes virtual card
	Remove(virtualCard *VirtualCard) error

	// Returns items owned by virtual card
	GetOwnedItems(virtualCard *VirtualCard) ([]OwnedItem, error)

	// Returns items owned by virtual card with PublicIds that match ids
	FilterOwnedItems(virtualCard *VirtualCard, ids []string) ([]OwnedItem, error)

	// Creates a new OwnedItem for virtual card if all conditions are met
	// (ex. virtualCard has enough points, item is still/already valid,
	// virtualCard has less items of this type than ItemDefinition.MaxAmount...)
	// itemDefinitionId passed as string - caller is not required to have "access" to the ItemDefinition object
	BuyItem(virtual *VirtualCard, itemDefinitionId string) (*OwnedItem, error)

	// Returns item - item is removed from virtualCard, points are returned to the card.
	// Fails if item was used.
	ReturnItem(ownedItem *OwnedItem) error
}

type VirtualCardManagerImpl struct {
	baseServices BaseServices
}

func CreateVirtualCardManagerImpl(baseServices BaseServices) *VirtualCardManagerImpl {
	return &VirtualCardManagerImpl{
		baseServices: baseServices,
	}
}

func (manager *VirtualCardManagerImpl) Create(user *User, businessId string) (*VirtualCard, error) {
	var virtualCard VirtualCard
	err := manager.baseServices.Database.Transaction(func(tx GormDB) error {
		// Find business by id
		var business Business
		businessResult := tx.First(&business, Business{PublicId: businessId})
		err := businessResult.GetError()
		if err == gorm.ErrRecordNotFound {
			return ErrNoSuchBusiness
		} else if err != nil {
			return fmt.Errorf("tx.First returned an error: %+v", err)
		}

		// Checks if user already has this virtual card
		// top 1 gorm pitfalls: do not query by relationship objects
		// or idk why code below just returns the first ever card
		// result := tx.First(&virtualCard, VirtualCard{User: user, Business: &business})
		result := tx.First(&virtualCard, VirtualCard{OwnerId: user.ID, BusinessId: business.ID})
		err = result.GetError()
		if err != gorm.ErrRecordNotFound && err != nil {
			return fmt.Errorf("tx.First returned an error: %+v", err)
		} else if err == nil {
			return ErrVirtualCardAlreadyExists
		}

		// Creates the card
		virtualCard = VirtualCard{
			PublicId: shortuuid.New(),
			Points:   0,
			User:     user,
			Business: &business,
		}

		result = tx.Create(&virtualCard)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("tx.Create returned an error: %+v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &virtualCard, nil
}

func (manager *VirtualCardManagerImpl) Remove(virtualCard *VirtualCard) error {
	//NOTE: this does not actually remove the rows from the database. All models contain gorm.Model
	//and that means, that soft deletes are done instead https://gorm.io/docs/delete.html#Soft-Delete
	//I'm not sure if it makes sense to remove transactions too. We could still allow businesses to see old
	//transactions by using Unscoped. I doubt that this is a very good design, but we don't have any
	//matching requirements
	return manager.baseServices.Database.Transaction(func(db GormDB) error {
		var tmp VirtualCard
		// Updates card data
		result := db.First(&tmp, &VirtualCard{PublicId: virtualCard.PublicId})
		err := result.GetError()
		if err == gorm.ErrRecordNotFound {
			return ErrNoSuchVirtualCard
		} else if err != nil {
			return fmt.Errorf("db.First returned an error: %+v", err)
		}

		// Deletes all items owned by card
		result = db.Where("virtual_card_id = ?", virtualCard.ID).Delete(&OwnedItem{})
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Delete(OwnedItem) returned an error: %+v", err)
		}

		// Deletes all TransactionDetails from transactions created with card
		//can't figure out how to delete with orm's built in methods
		//returns ErrMissingWhereClause no matter what i try
		result = db.Exec(`update transaction_details as td 
			set deleted_at = now()
			from transactions t
			where
				t.id = td.transaction_id and
				t.virtual_card_id = ?`, virtualCard.ID)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Delete(TransactionDetail) returned an error: %+v", err)
		}

		// Deletes all transactions owned by card
		result = db.Where("virtual_card_id = ?", virtualCard.ID).Delete(&Transaction{})
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Delete(Transaction) returned an error: %+v", err)
		}

		// Deletes the card
		result = db.Delete(virtualCard)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Delete(VirtualCard) returned an error: %+v", err)
		}

		return nil
	})
}

// idk if this is even necessary
// perhaps adding preloading functionality to accessors would work better
func (manager *VirtualCardManagerImpl) GetOwnedItems(virtualCard *VirtualCard) ([]OwnedItem, error) {
	var ownedItems []OwnedItem
	result := manager.baseServices.Database.Find(&ownedItems, &OwnedItem{VirtualCardId: virtualCard.ID})
	if err := result.GetError(); err != nil {
		return nil, err
	}
	return ownedItems, nil
}

func (manager *VirtualCardManagerImpl) FilterOwnedItems(virtualCard *VirtualCard, ids []string) ([]OwnedItem, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var ownedItems []OwnedItem
	result := manager.baseServices.Database.
		Where("virtual_card_id = ?", virtualCard.ID).
		Where("public_id in ?", ids).
		Find(&ownedItems)
	if err := result.GetError(); err != nil {
		return nil, err
	}
	return ownedItems, nil
}

// Verifies if virtualCard has less items of type itemDefinition than the allowed amount
func verifyMaxAmount(db GormDB, virtualCard *VirtualCard, itemDefinition *ItemDefinition) error {
	if itemDefinition.MaxAmount != 0 {
		var amount uint
		result := db.Model(&OwnedItem{}).
			Select("count(*)").
			Where(&OwnedItem{
				VirtualCard:    virtualCard,
				ItemDefinition: itemDefinition,
				Used:           sql.NullTime{Valid: false},
				Status:         OwnedItemStatusOwned,
			}).
			Scan(&amount)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Select(count(*)) returned an error %+v", err)
		}

		if amount >= itemDefinition.MaxAmount {
			return ErrAboveMaxAmount
		}
	}
	return nil
}

// Verifies if virtualCard has enough points to buy item of type itemDefinition
func verifyPrice(db GormDB, virtualCard *VirtualCard, itemDefinition *ItemDefinition) error {
	var points uint
	result := db.Model(&VirtualCard{}).
		Select("points").
		Where("id = ?", virtualCard.ID).
		Scan(&points)

	if err := result.GetError(); err != nil {
		return fmt.Errorf("db.Select(points) returned an error %+v", err)
	}

	if itemDefinition.Price > points {
		return ErrNotEnoughPoints
	}

	return nil
}

func (manager *VirtualCardManagerImpl) BuyItem(virtualCard *VirtualCard, itemDefinitionId string) (*OwnedItem, error) {
	var ownedItem OwnedItem
	err := manager.baseServices.Database.Transaction(func(db GormDB) error {
		// Find itemDefinition by id, handle errors
		var itemDefinition ItemDefinition
		itemDefinitionResult := db.First(&itemDefinition, &ItemDefinition{
			PublicId:   itemDefinitionId,
			BusinessId: virtualCard.BusinessId,
		})
		err := itemDefinitionResult.GetError()
		if err == gorm.ErrRecordNotFound {
			return ErrNoSuchItemDefinition
		} else if err != nil {
			return fmt.Errorf("tx.First returned an error: %+v", err)
		}

		// Updates itemDefinition
		result := db.First(&itemDefinition, "id = ?", itemDefinition.ID)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.First(itemDefinition) returned an error %+v", err)
		}

		// Checks if itemDefinition is valid
		if itemDefinition.Withdrawn {
			return ErrWithdrawnItem
		} else if !itemDefinition.Available {
			return ErrUnavailableItem
		} else if itemDefinition.StartDate.Valid && time.Now().Before(itemDefinition.StartDate.Time) {
			return ErrBeforeStartDate
		} else if itemDefinition.EndDate.Valid && time.Now().After(itemDefinition.EndDate.Time) {
			return ErrAfterEndDate
		} else if err := verifyMaxAmount(db, virtualCard, &itemDefinition); err != nil {
			return err
		} else if err := verifyPrice(db, virtualCard, &itemDefinition); err != nil {
			return err
		}

		// Creates the item
		ownedItem = OwnedItem{
			PublicId:       shortuuid.New(),
			Used:           sql.NullTime{Valid: false},
			Status:         OwnedItemStatusOwned,
			ItemDefinition: &itemDefinition,
			VirtualCard:    virtualCard,
		}

		// Subtracts points from card
		virtualCard.Points -= itemDefinition.Price
		result = db.Save(virtualCard)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Save(virtualCard) returned an error %+v", err)
		}

		// Creates ownedItem
		result = db.Create(&ownedItem)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.First(itemDefinition) returned an error %+v", err)
		}

		return nil
	}, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	return &ownedItem, nil
}

func (manager *VirtualCardManagerImpl) ReturnItem(ownedItem *OwnedItem) error {
	return manager.baseServices.Database.Transaction(func(db GormDB) error {
		// Loads VirtualCard and ItemDefinition relations in ownedItem
		result := db.
			Preload("VirtualCard").
			Preload("ItemDefinition").
			Find(ownedItem, "id = ?", ownedItem.ID)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Find(ownedItem) returned an error %+v", err)
		}

		// Checks if item is owned and was not used yet
		if ownedItem.Status != OwnedItemStatusOwned || ownedItem.Used.Valid {
			return ErrItemCantBeReturned
		}

		// Modifies item status
		ownedItem.Status = OwnedItemStatusReturned
		// Returns points to VirtualCard that owns ownedItem
		ownedItem.VirtualCard.Points += ownedItem.ItemDefinition.Price

		// Saves ownedItem
		result = db.Save(ownedItem)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Save(ownedItem) returned an error %+v", err)
		}

		// Saves the VirtualCard that owns OwnedItem
		//just to be sure
		result = db.Save(ownedItem.VirtualCard)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Save(ownedItem.VirtualCard) returned an error %+v", err)
		}

		return nil
	}, &sql.TxOptions{Isolation: sql.LevelSerializable})
}
