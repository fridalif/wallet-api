package repos_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockPool struct {
	mock.Mock
}

func (m *MockPool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	mockArgs := m.Called(ctx, sql, args)
	return mockArgs.Get(0).(pgconn.CommandTag), mockArgs.Error(1)
}

func (m *MockPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
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
			repo := &repos.wallet
		})
	}
}
