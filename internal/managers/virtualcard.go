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

var VirtualCardAlreadyExists = errors.New("Virtual card already exists")
var NoSuchVirtualCard = errors.New("No such virtual card")
var AboveMaxAmount = errors.New("Attempt to buy item above max item amount")
var NotEnoughPoints = errors.New("Attempt to buy item with not enough points")
var BeforeStartDate = errors.New("Attempt to buy item before start date")
var AfterEndDate = errors.New("Attempt to buy item after end date")
var UnavailableItem = errors.New("Attempt to buy an unavailable item")
var WithdrawnItem = errors.New("Attempt to buy a withdrawn item")
var ItemCantBeReturned = errors.New("Item can't be returned")

type VirtualCardManager interface {
	Create(user *User, business *Business) (*VirtualCard, error)
	Remove(virtualCard *VirtualCard) error
	GetOwnedItems(virtualCard *VirtualCard) ([]OwnedItem, error)
	FilterOwnedItems(virtualCard *VirtualCard, ids []string) ([]OwnedItem, error)
	BuyItem(virtual *VirtualCard, itemDefinition *ItemDefinition) (OwnedItem, error)
	ReturnItem(ownedItem *OwnedItem) error
}

type VirtualCardManagerImpl struct {
	baseServices *BaseServices
}

func (manager *VirtualCardManagerImpl) Create(user User, business Business) (*VirtualCard, error) {
	var virtualCard VirtualCard
	err := manager.baseServices.Database.Transaction(func(tx GormDB) error {
		result := tx.First(&virtualCard, VirtualCard{User: &user, Business: &business})
		err := result.GetError()
		if err != gorm.ErrRecordNotFound && err != nil {
			return fmt.Errorf("tx.First returned an error: %+v", err)
		} else if err == nil {
			return VirtualCardAlreadyExists
		}

		virtualCard = VirtualCard{
			PublicId: shortuuid.New(),
			Points:   0,
			User:     &user,
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
		result := db.First(&tmp, &VirtualCard{PublicId: virtualCard.PublicId})
		err := result.GetError()
		if err == gorm.ErrRecordNotFound {
			return NoSuchVirtualCard
		} else if err != nil {
			return fmt.Errorf("db.First returned an error: %+v", err)
		}

		result = db.Where("virtual_card_id = ?", virtualCard.ID).Delete(&OwnedItem{})
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Delete(OwnedItem) returned an error: %+v", err)
		}

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

		result = db.Where("virtual_card_id = ?", virtualCard.ID).Delete(&Transaction{})
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Delete(Transaction) returned an error: %+v", err)
		}

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
	var ownedItems []OwnedItem
	result := manager.baseServices.Database.Where("public_id in ?", ids).Find(&ownedItems)
	if err := result.GetError(); err != nil {
		return nil, err
	}
	return ownedItems, nil
}

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
			return AboveMaxAmount
		}
	}
	return nil
}

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
		return NotEnoughPoints
	}

	return nil
}

func (manager *VirtualCardManagerImpl) BuyItem(virtualCard *VirtualCard, itemDefinition *ItemDefinition) (*OwnedItem, error) {
	var ownedItem OwnedItem
	err := manager.baseServices.Database.Transaction(func(db GormDB) error {
		result := db.First(itemDefinition, "id = ?", itemDefinition.ID)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.First(itemDefinition) returned an error %+v", err)
		}

		if itemDefinition.Withdrawn {
			return WithdrawnItem
		} else if !itemDefinition.Available {
			return UnavailableItem
		} else if itemDefinition.StartDate.Valid && time.Now().Before(itemDefinition.StartDate.Time) {
			return BeforeStartDate
		} else if itemDefinition.EndDate.Valid && time.Now().After(itemDefinition.EndDate.Time) {
			return AfterEndDate
		} else if err := verifyMaxAmount(db, virtualCard, itemDefinition); err != nil {
			return err
		} else if err := verifyPrice(db, virtualCard, itemDefinition); err != nil {
			return err
		}

		ownedItem = OwnedItem{
			PublicId:       shortuuid.New(),
			Used:           sql.NullTime{Valid: false},
			Status:         OwnedItemStatusOwned,
			ItemDefinition: itemDefinition,
			VirtualCard:    virtualCard,
		}

		virtualCard.Points -= itemDefinition.Price
		result = db.Save(virtualCard)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Save(virtualCard) returned an error %+v", err)
		}

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
		result := db.
			Preload("VirtualCard").
			Preload("ItemDefinition").
			Find(ownedItem, "id = ?", ownedItem.ID)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Find(ownedItem) returned an error %+v", err)
		}

		if ownedItem.Status != OwnedItemStatusOwned || ownedItem.Used.Valid {
			return ItemCantBeReturned
		}

		ownedItem.Status = OwnedItemStatusReturned
		ownedItem.VirtualCard.Points += ownedItem.ItemDefinition.Price

		result = db.Save(ownedItem)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Save(ownedItem) returned an error %+v", err)
		}

		//just to be sure
		result = db.Save(ownedItem.VirtualCard)
		if err := result.GetError(); err != nil {
			return fmt.Errorf("db.Save(ownedItem.VirtualCard) returned an error %+v", err)
		}

		return nil
	}, &sql.TxOptions{Isolation: sql.LevelSerializable})
}
