package managers

import (
	"errors"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services"
	"github.com/lithammer/shortuuid/v4"
	"gorm.io/gorm"
)

var InvalidCardType = errors.New("Invalid card type")
var CardAlreadyExists = errors.New("Card already exists")
var CardDoesNotExist = errors.New("Card does not exist")

type LocalCardManager interface {
	Create(user *User, details *LocalCardDetails) (LocalCard, error)
	Remove(card *LocalCard) error
}

type LocalCardDetails struct {
	Type string
	Code string
	Name string
}

type LocalCardManagerImpl struct {
	baseServices *BaseServices
}

func (manager *LocalCardManagerImpl) Create(user *User, details *LocalCardDetails) (*LocalCard, error) {
	var cardType CardType
	for _, t := range CardTypes {
		if t.PublicId == details.Type {
			cardType = t
			break
		}
	}
	if cardType == (CardType{}) {
		return nil, InvalidCardType
	}

	//TODO verify code

	var foundLocalCard LocalCard
	tx := manager.baseServices.Database.First(&foundLocalCard, LocalCard{
		OwnerId: user.ID,
		Code:    details.Code,
		Type:    cardType.PublicId,
	})
	err := tx.GetError()
	if err == nil {
		return nil, CardAlreadyExists
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	localCard := LocalCard{
		PublicId: shortuuid.New(),
		Type:     cardType.PublicId,
		Name:     details.Name,
		Code:     details.Code,
		User:     user,
	}
	tx = manager.baseServices.Database.Create(&localCard)
	if err := tx.GetError(); err != nil {
		return nil, err
	}
	return &localCard, nil
}

func (manager *LocalCardManagerImpl) Remove(card *LocalCard) error {
	tx := manager.baseServices.Database.Delete(card)
	err := tx.GetError()
	if err == gorm.ErrRecordNotFound || tx.GetRowsAffected() == 0 {
		return CardDoesNotExist
	} else if err != nil {
		return err
	}
	return nil
}
