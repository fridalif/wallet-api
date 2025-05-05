package repos

import (
	"backend/pkg/wallet"

	"github.com/google/uuid"
)

type WalletRepositoryI interface {
	InitTable()
	GetWallet(id uuid.UUID) (*wallet.Wallet
	UpdateWallet()
}
