package services

import "github.com/google/uuid"

type WalletServiceI interface {
	GetBalance(id uuid.UUID) (int64, error)
	UpdateBalance(id uuid.UUID, operationType string, amount int64) error
}
