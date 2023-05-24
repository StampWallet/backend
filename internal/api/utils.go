package api

import (
	"fmt"

	"github.com/StampWallet/backend/internal/api/models"
	"github.com/StampWallet/backend/internal/database"
)

// on one hand, this is not necessary. both enums have the same values
// on the other, this could change, and I would prefer a crash if that ever happens
// it seems enums in go provide exactly zero type safety - just another type alias for string
func convertApiTransactionState(arg api.TransactionStateEnum) database.TransactionStateEnum {
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
