package repos

import (
	"backend/pkg/config"
	"backend/pkg/customerror"
	"backend/pkg/wallet"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepositoryI interface {
	CreateTables(ctx context.Context) error
	GetWallet(id uuid.UUID, ctx context.Context) (*wallet.Wallet, error)
	UpdateWallet(ctx context.Context) error
}

type walletRepository struct {
	Pool *pgxpool.Pool
	Host string
	Port string
}

func NewWalletRepository(appConfig *config.Config) (WalletRepositoryI, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", appConfig.DbUser, appConfig.DbPassword, appConfig.DbHost, appConfig.DbPort, appConfig.DbName)
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return &walletRepository{}, customerror.NewError("NewWalletRepository", appConfig.WebHost+":"+appConfig.WebPort, err.Error())
	}
	config.MaxConns = 100
	config.MinConns = 10
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 15 * time.Minute
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return &walletRepository{}, customerror.NewError("NewWalletRepository", appConfig.WebHost+":"+appConfig.WebPort, err.Error())
	}

	if err := pool.Ping(context.Background()); err != nil {
		return &walletRepository{}, customerror.NewError("NewWalletRepository", appConfig.WebHost+":"+appConfig.WebPort, err.Error())
	}
	return &walletRepository{
		Pool: pool,
		Host: appConfig.WebHost,
		Port: appConfig.WebPort,
	}, nil
}

func (walletRepo walletRepository) CreateTables(ctx context.Context) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS wallet (
		id UUID PRIMARY KEY,
		amount BIGINT,
	);`
	_, err := walletRepo.Pool.Exec(ctx, createTableQuery)
	if err != nil {
		return customerror.NewError("walletRepo.CreateTables", walletRepo.Host+":"+walletRepo.Port, err.Error())
	}
	createIndexQuery := `CREATE INDEX IF NOT EXISTS wallet_id_idx ON wallet(id);`
	_, err = walletRepo.Pool.Exec(ctx, createIndexQuery)
	if err != nil {
		return customerror.NewError("walletRepo.CreateTables", walletRepo.Host+":"+walletRepo.Port, err.Error())
	}
	return nil
}

func (walletRepo walletRepository) GetWallet(id uuid.UUID, ctx context.Context) (*wallet.Wallet, error) {
	var wallet wallet.Wallet
	selectQuery := "SELECT id, amount FROM wallet WHERE id = $1"
	err := walletRepo.Pool.QueryRow(ctx, selectQuery, id).Scan(&wallet.ID, &wallet.Amount)
	if err == nil {
		return &wallet, nil
	}
	if err == pgx.ErrNoRows {
		return nil, err
	}
	return nil, customerror.NewError("walletRepo.GetWallet", walletRepo.Host+":"+walletRepo.Port, err.Error())
}

func (walletRepo walletRepository) UpdateWallet(ctx context.Context) error {

}
