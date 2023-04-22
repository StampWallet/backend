package managers

import (
    "time"
    . "github.com/StampWallet/backend/internal/database"
    . "github.com/StampWallet/backend/internal/services"
)

type ItemDefinitionManager interface {
    AddItem(business *Business, details *ItemDetails) (*ItemDefinition, error)
    ChangeItemDetails(item *ItemDefinition, details *ItemDetails) (*ItemDefinition, error)
    WithdrawItem(item *ItemDefinition) error
    GetForBusiness(business *Business) ([]ItemDetails, error)
}

type ItemDetails struct {
    Name string
    Price *uint64
    Description string
    ImageId string
    StartDate *time.Time
    EndDate *time.Time
    MaxAmount *uint64
    Available *bool
}

type ItemDefinitionManagerImpl struct {
    baseServices *BaseServices
}

func (manager *ItemDefinitionManagerImpl) AddItem(business *Business, details *ItemDetails) (*ItemDefinition, error) {
    return nil, nil
}

func (manager *ItemDefinitionManagerImpl) ChangeItemDetails(item *ItemDefinition, details *ItemDetails) (*ItemDefinition, error) {
    return nil, nil
}

func (manager *ItemDefinitionManagerImpl) WithdrawItem(item *ItemDefinition) error {
    return nil
}

func (manager *ItemDefinitionManagerImpl) GetForBusiness(business *Business) ([]ItemDetails, error) {
    return nil, nil
}
