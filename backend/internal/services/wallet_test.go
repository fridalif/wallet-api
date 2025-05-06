package services_test

import (
	"backend/internal/services"
	"backend/pkg/customerror"
	"backend/pkg/wallet"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetWallet(ctx context.Context, id uuid.UUID) (*wallet.Wallet, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*wallet.Wallet), args.Error(1)
}

func (m *MockRepository) UpdateWallet(ctx context.Context, id uuid.UUID, delta int64) error {
	args := m.Called(ctx, id, delta)
	return args.Error(0)
}

func (m *MockRepository) CreateTables(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockRepository) ClosePull() {
	m.Called()
}

type GetBalanceTest struct {
	Name           string
	WalletId       uuid.UUID
	Mock           func(*MockRepository)
	WaitingBalance int64
	WaitingError   error
}

func TestWalletService_GetBalance(t *testing.T) {
	testID := uuid.New()
	testWallet := &wallet.Wallet{ID: testID, Amount: 100}

	tests := []GetBalanceTest{
		{
			Name:     "Success Test",
			WalletId: testID,
			Mock: func(r *MockRepository) {
				r.On("GetWallet", mock.Anything, testID).Return(testWallet, nil)
			},
			WaitingBalance: 100,
			WaitingError:   nil,
		},
		{
			Name:     "Not Found Test",
			WalletId: testID,
			Mock: func(r *MockRepository) {
				r.On("GetWallet", mock.Anything, testID).Return(&wallet.Wallet{}, pgx.ErrNoRows)
			},
			WaitingBalance: 0,
			WaitingError:   pgx.ErrNoRows,
		},
		{
			Name:     "Other Error Test",
			WalletId: testID,
			Mock: func(r *MockRepository) {
				r.On("GetWallet", mock.Anything, testID).Return(&wallet.Wallet{}, customerror.NewError("", "", "error"))
			},
			WaitingBalance: 0,
			WaitingError:   customerror.NewError("GetBalance.", "", "error"),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			test.Mock(mockRepo)

			service := services.NewWalletService(mockRepo)
			got, err := service.GetBalance(test.WalletId)
			if test.WaitingError != nil {
				assert.EqualError(t, err, test.WaitingError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.WaitingBalance, got)
			mockRepo.AssertExpectations(t)
		})
	}
}

type UpdateBalanceTest struct {
	Name          string
	WalletId      uuid.UUID
	OperationType string
	Amount        int64
	Mock          func(*MockRepository)
	WaitingError  error
}

func TestWalletService_UpdateBalance(t *testing.T) {
	testID := uuid.New()

	tests := []UpdateBalanceTest{
		{
			Name:          "Success Deposit Test",
			WalletId:      testID,
			OperationType: "DEPOSIT",
			Amount:        100,
			Mock: func(r *MockRepository) {
				r.On("UpdateWallet", mock.Anything, testID, int64(100)).Return(nil)
			},
			WaitingError: nil,
		},
		{
			Name:          "Success Withdraw Test",
			WalletId:      testID,
			OperationType: "WITHDRAW",
			Amount:        -100,
			Mock: func(r *MockRepository) {
				r.On("UpdateWallet", mock.Anything, testID, int64(-100)).Return(nil)
			},
			WaitingError: nil,
		},
		{
			Name:          "wrong operation",
			WalletId:      testID,
			OperationType: "INVALID",
			Amount:        100,
			Mock:          func(r *MockRepository) {},
			WaitingError:  customerror.ErrWrongOperation,
		},
		{
			Name:          "wrong amount",
			WalletId:      testID,
			OperationType: "WITHDRAW",
			Amount:        100,
			Mock: func(r *MockRepository) {
				r.On("UpdateWallet", mock.Anything, testID, int64(-1000)).Return(customerror.ErrWrongAmount)
			},
			WaitingError: customerror.ErrWrongAmount,
		},
		{
			Name:          "not found",
			WalletId:      testID,
			OperationType: "DEPOSIT",
			Amount:        100,
			Mock: func(r *MockRepository) {
				r.On("UpdateWallet", mock.Anything, testID, int64(100)).Return(pgx.ErrNoRows)
			},
			WaitingError: pgx.ErrNoRows,
		},
		{
			Name:          "other error",
			WalletId:      testID,
			OperationType: "DEPOSIT",
			Amount:        100,
			Mock: func(r *MockRepository) {
				r.On("UpdateWallet", mock.Anything, testID, int64(100)).Return(customerror.NewError("", "", "error"))
			},
			WaitingError: customerror.NewError("UpdateBalance.", "", "error"),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			test.Mock(mockRepo)

			service := services.NewWalletService(mockRepo)
			err := service.UpdateBalance(test.WalletId, test.OperationType, test.Amount)
			if test.WaitingError != nil {
				assert.EqualError(t, err, test.WaitingError.Error())
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
