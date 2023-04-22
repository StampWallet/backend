package managers

import (
    . "github.com/StampWallet/backend/internal/database"
    . "github.com/StampWallet/backend/internal/services"
)

//TODO
type ItemWithAction struct {
    PublicId string
    Action ActionTypeEnum
}

type TransactionManager interface {
    Start(card *VirtualCard, items []OwnedItem) (*Transaction, error)
    Finalize(transaction *Transaction, items []ItemWithAction, points uint64) (*Transaction, error)
}

type TransactionManagerImpl struct {
    baseServices *BaseServices
}

func (manager *TransactionManagerImpl) Start(card *VirtualCard, items []OwnedItem) (*Transaction, error) {
    return nil, nil
}

func (manager *TransactionManagerImpl) Finalize(transaction *Transaction, items []ItemWithAction, points uint64) (*Transaction, error) {
    return nil, nil
}
