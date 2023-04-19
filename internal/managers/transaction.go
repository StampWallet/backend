package manager

import (
    . "github.com/StampWallet/backend/internal/database"
)

type TransactionManager interface {
    Start(card *VirtualCard, items []OwnedItem) (*Transaction, error)
    Finalize(transaction *Transaction, items []ItemWithStatus, points uint64) (*Transaction, error)
}

type TransactionManagerImpl struct {
    baseServices *BaseServices
}

func (manager *TransactionManagerImpl) Start(card *VirtualCard, items []OwnedItem) (*Transaction, error) {

}

func (manager *TransactionManagerImpl) Finalize(transaction *Transaction, items []ItemWithStatus, points uint64) (*Transaction, error) {

}
