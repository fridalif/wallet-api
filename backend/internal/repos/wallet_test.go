package repos_test

import (
	"backend/internal/repos"
	"backend/pkg/customerror"
	"backend/pkg/wallet"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

type MockPool struct {
	mock.Mock
}

func (m *MockPool) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	mockArgs := m.Called(ctx, sql, args)
	return mockArgs.Get(0).(pgconn.CommandTag), mockArgs.Error(1)
}

func (m *MockPool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	mockArgs := m.Called(ctx, sql, args)
	return mockArgs.Get(0).(pgx.Row)
}

func (m *MockPool) Close() {
	m.Called()
}

type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...any) error {
	args := m.Called(dest)
	return args.Error(0)
}

type CreateTablesTest struct {
	Name    string
	Mock    func(*MockPool)
	WantErr bool
}

func TestWalletRepository_CreateTables(t *testing.T) {
	createTableTests := []CreateTablesTest{
		{
			Name: "Success Test",
			Mock: func(m *MockPool) {
				m.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, nil).Twice()
			},
			WantErr: false,
		},
		{
			Name: "Error creating table",
			Mock: func(m *MockPool) {
				m.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, errors.New("error creating table")).Once()
			},
			WantErr: true,
		},
		{
			Name: "Error creating index",
			Mock: func(m *MockPool) {
				m.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, errors.New("error creating index")).Once()
			},
			WantErr: true,
		},
	}

	for _, test := range createTableTests {
		t.Run(test.Name, func(t *testing.T) {
			mockPool := new(MockPool)
			test.Mock(mockPool)
			repo := &repos.WalletRepository{
				Pool: mockPool,
				Host: "localhost",
				Port: "80",
			}
			err := repo.CreateTables(context.Background())
			if test.WantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockPool.AssertExpectations(t)
		})
	}
}

type GetWalletTest struct {
	Name          string
	WalletId      uuid.UUID
	WaitingWallet *wallet.Wallet
	Mock          func(*MockPool, *MockRow)
	WaitingError  error
}

func TestWalletRepository_GetWallet(t *testing.T) {
	testUUID := uuid.New()
	testWallet := &wallet.Wallet{
		ID:     testUUID,
		Amount: 1000,
	}
	getWalletTests := []GetWalletTest{
		{
			Name:          "Success Test",
			WaitingError:  nil,
			WalletId:      testUUID,
			WaitingWallet: testWallet,
			Mock: func(p *MockPool, r *MockRow) {
				p.On("QueryRow", mock.Anything, "SELECT id, amount FROM wallet WHERE id = $1", []interface{}{testUUID}).Return(r)
				r.On("Scan", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					mockArgs := args.Get(0).([]interface{})
					idPtr := mockArgs[0].(*uuid.UUID)
					*idPtr = testWallet.ID
					amountPtr := mockArgs[1].(*int64)
					*amountPtr = testWallet.Amount
				}).Return(nil)
			},
		},
		{
			Name:          "Not Found Test",
			WaitingError:  pgx.ErrNoRows,
			WalletId:      testUUID,
			WaitingWallet: nil,
			Mock: func(p *MockPool, r *MockRow) {
				p.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(r)
				r.On("Scan", mock.Anything, mock.Anything).Return(pgx.ErrNoRows)
			},
		},
		{
			Name:          "Other Error Test",
			WaitingError:  customerror.NewError("walletRepo.GetWallet", "127.0.0.1:8080", errors.New("Other error").Error()),
			WalletId:      testUUID,
			WaitingWallet: nil,
			Mock: func(p *MockPool, r *MockRow) {
				p.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(r)
				r.On("Scan", mock.Anything, mock.Anything).Return(errors.New("Other error"))
			},
		},
	}
	for _, test := range getWalletTests {
		t.Run(test.Name, func(t *testing.T) {
			mockPool := new(MockPool)
			mockRow := new(MockRow)
			test.Mock(mockPool, mockRow)
			repo := &repos.WalletRepository{
				Pool: mockPool,
				Host: "127.0.0.1",
				Port: "8080",
			}
			gettedWallet, err := repo.GetWallet(context.Background(), test.WalletId)
			if err != nil {
				assert.ErrorIs(t, err, test.WaitingError)
				assert.Equal(t, gettedWallet, test.WaitingWallet)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, gettedWallet.ID, test.WaitingWallet.ID)
				assert.Equal(t, gettedWallet.Amount, test.WaitingWallet.Amount)
			}
			mockPool.AssertExpectations(t)
			mockRow.AssertExpectations(t)
		})
	}
}

/*
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
*/
type UpdateWalletTest struct {
	Name         string
	WalletId     uuid.UUID
	Delta        int64
	Mock         func(*MockPool)
	WaitingError error
}

func TestWalletRepository_UpdateWallet(t *testing.T) {
	testUUID := uuid.New()
	testDelta := int64(100)
	updateWalletTests := []UpdateWalletTest{
		{
			Name:         "Success Test",
			WalletId:     testUUID,
			WaitingError: nil,
			Delta:        testDelta,
			Mock: func(p *MockPool) {
				p.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pgconn.NewCommandTag("UPDATE 1"), nil)
			},
		},
		{
			Name:         "Not Found Test",
			WalletId:     testUUID,
			Delta:        testDelta,
			WaitingError: pgx.ErrNoRows,
			Mock: func(p *MockPool) {
				p.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pgconn.NewCommandTag("UPDATE 0"), nil)
			},
		},
		{
			Name:     "Wrong Amount Test",
			WalletId: testUUID,
			Delta:    testDelta,
			Mock: func(p *MockPool) {
				p.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, &pgconn.PgError{Code: "23514"})
			},
			WaitingError: customerror.ErrWrongAmount,
		},
		{
			Name:     "other error",
			WalletId: testUUID,
			Delta:    testDelta,
			Mock: func(p *MockPool) {
				p.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, errors.New("error"))
			},
			WaitingError: customerror.NewError("walletRepo.UpdateWallet", "127.0.0.1:8080", "error"),
		},
	}
	for _, test := range updateWalletTests {
		t.Run(test.Name, func(t *testing.T) {
			mockPool := new(MockPool)
			test.Mock(mockPool)

			repo := &repos.WalletRepository{
				Pool: mockPool,
				Host: "127.0.0.1",
				Port: "8080",
			}

			err := repo.UpdateWallet(context.Background(), test.WalletId, test.Delta)
			if test.WaitingError != nil {
				assert.EqualError(t, err, test.WaitingError.Error())
			} else {
				assert.NoError(t, err)
			}
			mockPool.AssertExpectations(t)
		})
	}
}
