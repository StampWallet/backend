package managers

import (
	"errors"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services"
)

var InvalidCardType = errors.New("Invalid card type")
var CardDoesNotExist = errors.New("Invalid card type")

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

func (handler *LocalCardManagerImpl) Create(user *User, details *LocalCardDetails) (*LocalCard, error) {
	return nil, nil
}

func (handler *LocalCardManagerImpl) Remove(card *LocalCard) error {
	return nil
}

func (handler *LocalCardManagerImpl) GetForUser(user *User) ([]LocalCard, error) {
	return nil, nil
}
