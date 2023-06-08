package apiUtils

import (
	"fmt"

	"github.com/StampWallet/backend/internal/api/models"
	"github.com/StampWallet/backend/internal/database"
)

// on one hand, this is not necessary. both enums have the same values
// on the other, this could change, and I would prefer a crash if that ever happens
// it seems enums in go provide exactly zero type safety - just another type alias for string
func ConvertApiTransactionState(arg api.TransactionStateEnum) database.TransactionStateEnum {
	if arg == api.STARTED {
		return database.TransactionStateStarted
	} else if arg == api.PROCESSING {
		return database.TransactionStateProcesing
	} else if arg == api.FINISHED {
		return database.TransactionStateFinished
	} else if arg == api.EXPIRED {
		return database.TransactionStateExpired
	} else {
		panic(fmt.Errorf("unkown api.TransactionStateEnum enum valule - cannot map to database.TransactionStateEnum %+v", arg))
	}
}

func ConvertDbTransactionState(arg database.TransactionStateEnum) api.TransactionStateEnum {
	if arg == database.TransactionStateStarted {
		return api.STARTED
	} else if arg == database.TransactionStateProcesing {
		return api.PROCESSING
	} else if arg == database.TransactionStateFinished {
		return api.FINISHED
	} else if arg == database.TransactionStateExpired {
		return api.EXPIRED
	} else {
		panic(fmt.Errorf("unkown api.TransactionStateEnum enum valule - cannot map to database.TransactionStateEnum %+v", arg))
	}
}

func ConvertApiItemAction(arg api.ItemActionTypeEnum) database.ActionTypeEnum {
	if arg == api.REDEEMED {
		return database.RedeemedActionType
	} else if arg == api.RECALLED {
		return database.RecalledActionType
	} else if arg == api.CANCELLED {
		return database.CancelledActionType
	} else {
		panic(fmt.Errorf("unkown api.TransactionStateEnum enum valule - cannot map to database.TransactionStateEnum %+v", arg))
	}
}
