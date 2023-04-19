package manager

import (
    . "github.com/StampWallet/backend/internal/database"
)

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

}

func (manager *VirtualCardManagerImpl) Remove(virtualCard *VirtualCard) error {

}

func (manager *VirtualCardManagerImpl) GetForUser(user *User) ([]VirtualCard, error) {

}

func (manager *VirtualCardManagerImpl) GetOwnedItems(virtualCard *VirtualCard) ([]OwnedItem, error) {

}

func (manager *VirtualCardManagerImpl) FilterOwnedItems(virtualCard *VirtualCard, ids []string) ([]OwnedItem, error) {

}

func (manager *VirtualCardManagerImpl) BuyItem(virtual *VirtualCard, itemDefinition *ItemDefinition) (*OwnedItem, error) {

}

func (manager *VirtualCardManagerImpl) ReturnItem(ownedItem *OwnedItem) error {

}

