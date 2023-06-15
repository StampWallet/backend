package managers

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services"
	"github.com/StampWallet/backend/internal/utils"
	"github.com/lithammer/shortuuid/v4"
)

const transactionCodeLength = 12

var (
	ErrInvalidItem        = errors.New("Invalid item")        // no such item or item already used
	ErrInvalidTransaction = errors.New("Invalid transaction") // transaction finished
	ErrItemBadCardId      = errors.New("Owned item vcard id does not match that of provided vcard")
	ErrInvalidAction      = errors.New("NoActionType is not a valid action when finalizing transaction")
	ErrInvalidActionSet   = errors.New("Invalid action set - does not match started transaction details")
)

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

func digitsSlice() []int {
	return []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
}

func generateCode() string {
	randomIntSlice := utils.RandomSlice(transactionCodeLength, digitsSlice())
	randomStrDigits := utils.Map(randomIntSlice, strconv.Itoa)
	randomDigitStr := strings.Join(randomStrDigits, "")
	return randomDigitStr
}

func CreateTransactionManagerImpl(baseServices BaseServices) *TransactionManagerImpl {
	return &TransactionManagerImpl{
		baseServices: baseServices,
	}
}

func (manager *TransactionManagerImpl) Start(card *VirtualCard, items []OwnedItem) (*Transaction, error) {
	for _, chosenItem := range items {
		if chosenItem.VirtualCardId != card.ID {
			return nil, ErrItemBadCardId
		}
		if chosenItem.Status != OwnedItemStatusOwned {
			return nil, ErrInvalidItem
		}
	}

	var transaction *Transaction
	err := manager.baseServices.Database.Transaction(func(tx GormDB) error {
		// __jm__ TODO handle code collisions
		transaction = &Transaction{
			PublicId:      shortuuid.New(),
			VirtualCardId: card.ID,
			Code:          generateCode(),
			State:         TransactionStateStarted,
			AddedPoints:   0,
		}
		res := tx.Create(transaction)
		if err := res.GetError(); err != nil {
			return err
		}

		transactionDetails := []TransactionDetail{}
		for _, chosenItem := range items {
			transactionDetails = append(transactionDetails, TransactionDetail{
				TransactionId: transaction.ID,
				ItemId:        chosenItem.ID,
				Action:        NoActionType,
			})
		}
		res = tx.Create(transactionDetails)
		if err := res.GetError(); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (manager *TransactionManagerImpl) Finalize(transaction *Transaction, actions []ItemWithAction, points uint64) (*Transaction, error) {
	failTransaction := false
	for _, chosenItem := range actions {
		if chosenItem.Item.VirtualCardId != transaction.VirtualCardId {
			return nil, ErrItemBadCardId
		}
		if chosenItem.Action == NoActionType {
			return nil, ErrInvalidAction
		}
	}

	itemIdToAction := make(map[uint]ActionTypeEnum)
	for _, itemWithAction := range actions {
		item := itemWithAction.Item
		action := itemWithAction.Action
		itemIdToAction[item.ID] = action
	}

	err := manager.baseServices.Database.Transaction(func(tx GormDB) error {
		result := tx.
			Preload("TransactionDetails").
			Preload("TransactionDetails.OwnedItem").
			Preload("TransactionDetails.OwnedItem.ItemDefinition").
			Preload("VirtualCard").
			Find(transaction, "id = ?", transaction.ID)
		if err := result.GetError(); err != nil {
			return err
		}

		tds := transaction.TransactionDetails
		if len(actions) != len(tds) {
			// TODO: log diff between item sets
			return ErrInvalidActionSet
		}

		for i := range tds {
			td := &tds[i] // need to modify tds slice to save after loop
			// item changed between transactions
			if td.OwnedItem.Status != OwnedItemStatusOwned {
				failTransaction = true
				return ErrInvalidItem
			}

			action, ok := itemIdToAction[td.OwnedItem.ID]
			if !ok {
				return ErrInvalidItem
			}

			td.Action = action
			switch action {
			case RedeemedActionType:
				td.OwnedItem.Used = sql.NullTime{Time: time.Now(), Valid: true}
				td.OwnedItem.Status = OwnedItemStatusUsed
			case RecalledActionType:
				td.OwnedItem.Status = OwnedItemStatusWithdrawn
				transaction.VirtualCard.Points += td.OwnedItem.ItemDefinition.Price
			case CancelledActionType:
				// ?
			}

			result = tx.Save(td.OwnedItem)
			if err := result.GetError(); err != nil {
				return err
			}
		}
		result = tx.Save(tds)
		if err := result.GetError(); err != nil {
			return err
		}

		transaction.State = TransactionStateFinished
		transaction.AddedPoints = uint(points)
		transaction.VirtualCard.Points += transaction.AddedPoints

		result = tx.Save(transaction.VirtualCard)
		if err := result.GetError(); err != nil {
			return err
		}

		result = tx.Save(transaction)
		if err := result.GetError(); err != nil {
			return err
		}

		return nil
	})

	if failTransaction {
		transaction.State = TransactionStateFailed
		result := manager.baseServices.Database.Save(&transaction)
		if err := result.GetError(); err != nil {
			return nil, err
		}
	}

	return transaction, err
}
