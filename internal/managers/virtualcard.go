package managers

import (
	"errors"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services"
)

var VirtualCardAlreadyExists = errors.New("Virtual card already exists")
var NoSuchVirtualCard = errors.New("No such virtual card")
var AboveMaxAmount = errors.New("Attempt to buy item above max item amount")
var NotEnoughPoints = errors.New("Attempt to buy item with not enough points")
var BeforeStartDate = errors.New("Attempt to buy item before start date")
var AfterEndDate = errors.New("Attempt to buy item after end date")
var UnavailableItem = errors.New("Attempt to buy an unavailable item")
var WithdrawnItem = errors.New("Attempt to buy a withdrawn item")
var ItemReturned = errors.New("Item already returned")

type VirtualCardManager interface {
	Create(user *User, business *Business) (*VirtualCard, error)
	Remove(virtualCard *VirtualCard) error
	GetForUser(user *User) ([]VirtualCard, error)
	GetOwnedItems(virtualCard *VirtualCard) ([]OwnedItem, error)
	FilterOwnedItems(virtualCard *VirtualCard, ids []string) ([]OwnedItem, error)
	BuyItem(virtual *VirtualCard, itemDefinition *ItemDefinition) (OwnedItem, error)
	ReturnItem(ownedItem *OwnedItem) error
}

type VirtualCardManagerImpl struct {
	baseServices *BaseServices
}

func (manager *VirtualCardManagerImpl) Create(user User, business Business) (*VirtualCard, error) {
	return nil, nil
}

func (manager *VirtualCardManagerImpl) Remove(virtualCard *VirtualCard) error {
	return nil
}

func (manager *VirtualCardManagerImpl) GetForUser(user *User) ([]VirtualCard, error) {
	return nil, nil
}

func (manager *VirtualCardManagerImpl) GetOwnedItems(virtualCard *VirtualCard) ([]OwnedItem, error) {
	return nil, nil
}

func (manager *VirtualCardManagerImpl) FilterOwnedItems(virtualCard *VirtualCard, ids []string) ([]OwnedItem, error) {
	return nil, nil
}

func (manager *VirtualCardManagerImpl) BuyItem(virtual *VirtualCard, itemDefinition *ItemDefinition) (*OwnedItem, error) {
	return nil, nil
}

func (manager *VirtualCardManagerImpl) ReturnItem(ownedItem *OwnedItem) error {
	return nil
}
