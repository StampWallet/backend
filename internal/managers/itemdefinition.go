package managers

import (
	"time"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services"
)

type ItemDefinitionManager interface {
	AddItem(business *Business, details *ItemDetails) (*ItemDefinition, error)
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

func (manager *ItemDefinitionManagerImpl) AddItem(business *Business, details *ItemDetails) (*ItemDefinition, error) {
	return nil, nil
}

func (manager *ItemDefinitionManagerImpl) ChangeItemDetails(item *ItemDefinition, details *ItemDetails) (*ItemDefinition, error) {
	return nil, nil
}

func (manager *ItemDefinitionManagerImpl) WithdrawItem(item *ItemDefinition) (*ItemDefinition, error) {
	return nil, nil
}

func (manager *ItemDefinitionManagerImpl) GetForBusiness(business *Business) ([]ItemDefinition, error) {
	return nil, nil
}
