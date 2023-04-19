package manager

import (
    . "github.com/StampWallet/backend/internal/database"
)

type LocalCardManager interface { 
    Create(user *User, details *LocalCardDetails) (LocalCard, error)
    Remove(card *LocalCard) error
    GetForUser(user *User) ([]LocalCard, error)
}

type LocalCardDetails struct {
    Type string
    Code string
    Name string
}

type LocalCardManagerImpl struct {
    baseServices *BaseServices
}

func (handler *LocalCardManagerImpl) Create(user *User, details *LocalCardDetails) (LocalCard, error) {

}

func (handler *LocalCardManagerImpl) Remove(card *LocalCard) error {

}

func (handler *LocalCardManagerImpl) GetForUser(user *User) ([]LocalCard, error) {

}
