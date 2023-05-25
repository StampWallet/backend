package managers

import (
	"errors"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services"
)

var InvalidItem = errors.New("Invalid item")

// TODO
type ItemWithAction struct {
	Item   *OwnedItem
	Action ActionTypeEnum
}

type TransactionManager interface {
	Start(card *VirtualCard, items []OwnedItem) (*Transaction, error)
	Finalize(transaction *Transaction, items []ItemWithAction, points uint64) (*Transaction, error)
}

type TransactionManagerImpl struct {
	baseServices BaseServices
}

func CreateTransactionManagerImpl(baseServices BaseServices) *TransactionManagerImpl {
	return &TransactionManagerImpl{
		baseServices: baseServices,
	}
}

func (manager *TransactionManagerImpl) Start(card *VirtualCard, items []OwnedItem) (*Transaction, error) {
	return nil, nil
}

func (manager *TransactionManagerImpl) Finalize(transaction *Transaction, items []ItemWithAction, points uint64) (*Transaction, error) {
	return nil, nil
}
