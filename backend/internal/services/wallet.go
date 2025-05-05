package services

import (
	"backend/internal/repos"
	"backend/pkg/customerror"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type WalletServiceI interface {
	GetBalance(id uuid.UUID) (int64, error)
	UpdateBalance(id uuid.UUID, operationType string, amount int64) error
}

type WalletService struct {
	Repo repos.WalletRepositoryI
}

func NewWalletService(repo repos.WalletRepositoryI) WalletServiceI {
	return &WalletService{
		Repo: repo,
	}
}

func (WalletService *WalletService) GetBalance(id uuid.UUID) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	wallet, err := WalletService.Repo.GetWallet(ctx, id)
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
func (WalletService *WalletService) UpdateBalance(id uuid.UUID, operationType string, amount int64) error {
	if operationType != "DEPOSIT" && operationType != "WITHDRAW" {
		return customerror.ErrWrongOperation
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := WalletService.Repo.UpdateWallet(ctx, id, amount)
	if err == nil || err == customerror.ErrWrongAmount {
		return err
	}
	customError := err.(customerror.CustomError)
	customError.AppendModule("UpdateBalance")
	return customError
}
