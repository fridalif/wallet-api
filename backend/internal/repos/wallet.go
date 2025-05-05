package repos

import (
	"backend/pkg/customerror"
	"backend/pkg/wallet"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepositoryI interface {
	InitTable()
	GetWallet(id uuid.UUID) (*wallet.Wallet, *customerror.CustomError)
	UpdateWallet() *customerror.CustomError
}

type walletRepository struct {
	Pool *pgxpool.Pool
}
