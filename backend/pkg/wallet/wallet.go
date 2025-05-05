package wallet

import "github.com/google/uuid"

type Wallet struct {
	ID     uuid.UUID
	Amount int64
}
