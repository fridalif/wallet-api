package requests

import "github.com/google/uuid"

type UpdateBalanceRequest struct {
	WalletId      uuid.UUID `json:"valletId"`
	OperationType string    `json:"operationType"`
	Amount        int64     `json:"amount"`
}
