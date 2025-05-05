package repos

import (
	"backend/pkg/config"
	"backend/pkg/customerror"
	"backend/pkg/wallet"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
)

type WalletRepositoryI interface {
	CreateTables() error
	GetWallet(id uuid.UUID) (*wallet.Wallet, error)
	UpdateWallet() error
}

type walletRepository struct {
	Pool *pgxpool.Pool
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
	}, nil
}

func (walletRepo walletRepository) CreateTables() error {

}
func (walletRepo walletRepository) GetWallet(id uuid.UUID) (*wallet.Wallet, error) {

}

func (walletRepo walletRepository) UpdateWallet() error {

}
