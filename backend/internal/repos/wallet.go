package repos

import (
	"backend/pkg/customerror"
	"backend/pkg/wallet"

	"github.com/google/uuid"
)

type WalletRepositoryI interface {
	InitTable()
	GetWallet(id uuid.UUID) (*wallet.Wallet, *customerror.CustomError)
	UpdateWallet() *customerror.CustomError
}
