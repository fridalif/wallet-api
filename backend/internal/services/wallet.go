package services

import (
	"backend/internal/repos"
	"backend/pkg/customerror"
	"backend/pkg/wallet"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type WalletServiceI interface {
	GetBalance(id uuid.UUID) (int64, error)
	UpdateBalance(id uuid.UUID, operationType string, amount int64) error
}

type walletService struct {
	Repo repos.WalletRepositoryI
}

func NewWalletService(repo repos.WalletRepositoryI) WalletServiceI {
	return &walletService{
		Repo: repo,
	}
}

func (walletService walletService) GetBalance(id uuid.UUID) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	wallet, err := walletService.Repo.GetWallet(ctx, id)
	if err == nil {
		return wallet.Amount, nil
	}
	if err == pgx.ErrNoRows {
		return 0, err
	}
	customError := err.(customerror.CustomError)
	customError.AppendModule("GetBalance")
	return 0, customError
}
func (walletService walletService) UpdateBalance(id uuid.UUID, operationType string, amount int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	wallet := wallet.Wallet{
		ID: id,
		Amount: ,
	}
}
