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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepositoryI interface {
	CreateTables(ctx context.Context) error
	GetWallet(ctx context.Context, id uuid.UUID) (*wallet.Wallet, error)
	UpdateWallet(ctx context.Context, id uuid.UUID, delta int64) error
	ClosePull()
}

type PoolInterface interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Close()
}

type WalletRepository struct {
	Pool PoolInterface
	Host string
	Port string
}

func NewWalletRepository(appConfig *config.Config) (WalletRepositoryI, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", appConfig.DbUser, appConfig.DbPassword, appConfig.DbHost, appConfig.DbPort, appConfig.DbName)
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return &WalletRepository{}, customerror.NewError("NewWalletRepository", appConfig.WebHost+":"+appConfig.WebPort, err.Error())
	}
	config.MaxConns = 100
	config.MinConns = 10
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 15 * time.Minute
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return &WalletRepository{}, customerror.NewError("NewWalletRepository", appConfig.WebHost+":"+appConfig.WebPort, err.Error())
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return &WalletRepository{}, customerror.NewError("NewWalletRepository", appConfig.WebHost+":"+appConfig.WebPort, err.Error())
	}
	return &WalletRepository{
		Pool: pool,
		Host: appConfig.WebHost,
		Port: appConfig.WebPort,
	}, nil
}

func (walletRepo *WalletRepository) CreateTables(ctx context.Context) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS wallet (
		id UUID PRIMARY KEY,
		amount BIGINT NOT NULL DEFAULT 0 CHECK (amount >= 0)
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

func (walletRepo *WalletRepository) GetWallet(ctx context.Context, id uuid.UUID) (*wallet.Wallet, error) {
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

func (walletRepo *WalletRepository) UpdateWallet(ctx context.Context, id uuid.UUID, delta int64) error {
	updateQuery := "UPDATE wallet set amount = amount + $1 WHERE id = $2"
	command, err := walletRepo.Pool.Exec(ctx, updateQuery, delta, id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23514" {
				return customerror.ErrWrongAmount
			}
		}
		return customerror.NewError("walletRepo.UpdateWallet", walletRepo.Host+":"+walletRepo.Port, err.Error())
	}
	if command.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (walletRepo *WalletRepository) ClosePull() {
	walletRepo.Pool.Close()
}
