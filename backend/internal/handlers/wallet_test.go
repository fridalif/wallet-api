package handlers_test

import (
	"backend/internal/handlers"
	"backend/pkg/customerror"
	"backend/pkg/requests"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetBalance(id uuid.UUID) (int64, error) {
	args := m.Called(id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockService) UpdateBalance(id uuid.UUID, operationType string, amount int64) error {
	args := m.Called(id, operationType, amount)
	return args.Error(0)
}

type GetBalanceTest struct {
	Name           string
	WalletId       string
	Mock           func(*MockService)
	ExpectedStatus int
	ExpectedBody   gin.H
}

func TestWalletHandler_GetBalance(t *testing.T) {
	testID := uuid.New()

	tests := []GetBalanceTest{
		{
			Name:     "Success Test",
			WalletId: testID.String(),
			Mock: func(s *MockService) {
				s.On("GetBalance", testID).Return(int64(100), nil)
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: gin.H{
				"status": float64(200),
				"data": map[string]interface{}{
					"balance": float64(100),
				},
				"error": nil,
			},
		},
		{
			Name:           "Invalid UUID Test",
			WalletId:       "invalid",
			Mock:           func(s *MockService) {},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: gin.H{
				"status": float64(400),
				"data":   map[string]interface{}{},
				"error":  "Wrong uuid",
			},
		},
		{
			Name:     "Not Found Test",
			WalletId: testID.String(),
			Mock: func(s *MockService) {
				s.On("GetBalance", testID).Return(int64(0), pgx.ErrNoRows)
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: gin.H{
				"status": float64(404),
				"data":   map[string]interface{}{},
				"error":  "Wallet not found",
			},
		},
		{
			Name:     "Internal Server Error Test",
			WalletId: testID.String(),
			Mock: func(s *MockService) {
				s.On("GetBalance", testID).Return(int64(0), customerror.NewError("", "", "error"))
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: gin.H{
				"status": float64(500),
				"data":   map[string]interface{}{},
				"error":  "Internal Server Error",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			mockService := new(MockService)
			test.Mock(mockService)

			handler := handlers.NewWalletHandler(mockService)

			router := gin.Default()
			router.GET("/wallets/:id", handler.GetBalance)

			req, _ := http.NewRequest(http.MethodGet, "/wallets/"+test.WalletId, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, test.ExpectedStatus, resp.Code)

			var body gin.H
			err := json.Unmarshal(resp.Body.Bytes(), &body)
			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedBody, body)

			mockService.AssertExpectations(t)
		})
	}
}

type UpdateBalanceTest struct {
	Name           string
	Request        requests.UpdateBalanceRequest
	Mock           func(*MockService)
	ExpectedStatus int
	ExpectedBody   gin.H
}

func TestWalletHandler_UpdateBalance(t *testing.T) {
	testID := uuid.New()

	tests := []UpdateBalanceTest{
		{
			Name: "Success Deposit Test",
			Request: requests.UpdateBalanceRequest{
				WalletId:      testID,
				OperationType: "DEPOSIT",
				Amount:        100,
			},
			Mock: func(s *MockService) {
				s.On("UpdateBalance", testID, "DEPOSIT", int64(100)).Return(nil)
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: gin.H{
				"status": float64(200),
				"body":   map[string]interface{}{},
				"error":  nil,
			},
		},
		{
			Name: "Invalid Operation Type Test",
			Request: requests.UpdateBalanceRequest{
				WalletId:      testID,
				OperationType: "INVALID",
				Amount:        100,
			},
			Mock: func(s *MockService) {
				s.On("UpdateBalance", testID, "INVALID", int64(100)).Return(customerror.ErrWrongOperation)
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: gin.H{
				"status": float64(400),
				"body":   map[string]interface{}{},
				"error":  "Operation must be DEPOSIT or WITHDRAW",
			},
		},
		{
			Name: "Wrong Amount Test",
			Request: requests.UpdateBalanceRequest{
				WalletId:      testID,
				OperationType: "WITHDRAW",
				Amount:        1000,
			},
			Mock: func(s *MockService) {
				s.On("UpdateBalance", testID, "WITHDRAW", int64(1000)).Return(customerror.ErrWrongAmount)
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: gin.H{
				"status": float64(400),
				"body":   map[string]interface{}{},
				"error":  "Amount cant be less than zero",
			},
		},
		{
			Name: "Not Found Test",
			Request: requests.UpdateBalanceRequest{
				WalletId:      testID,
				OperationType: "DEPOSIT",
				Amount:        100,
			},
			Mock: func(s *MockService) {
				s.On("UpdateBalance", testID, "DEPOSIT", int64(100)).Return(pgx.ErrNoRows)
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: gin.H{
				"status": float64(404),
				"body":   map[string]interface{}{},
				"error":  "Wallet not found",
			},
		},
		{
			Name: "Internal Server Error Test",
			Request: requests.UpdateBalanceRequest{
				WalletId:      testID,
				OperationType: "DEPOSIT",
				Amount:        100,
			},
			Mock: func(s *MockService) {
				s.On("UpdateBalance", testID, "DEPOSIT", int64(100)).Return(customerror.NewError("", "", "error"))
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: gin.H{
				"status": float64(500),
				"data":   map[string]interface{}{},
				"error":  "Internal Server Error",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			mockService := new(MockService)
			test.Mock(mockService)

			handler := handlers.NewWalletHandler(mockService)

			router := gin.Default()
			router.POST("/wallet", handler.UpdateBalance)

			body, _ := json.Marshal(test.Request)
			req, _ := http.NewRequest(http.MethodPost, "/wallet", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, test.ExpectedStatus, resp.Code)

			var responseBody gin.H
			err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
			assert.NoError(t, err)
			fmt.Printf("%v", responseBody)
			assert.Equal(t, test.ExpectedBody["status"], responseBody["status"])
			assert.Equal(t, test.ExpectedBody["error"], responseBody["error"])

			mockService.AssertExpectations(t)
		})
	}
}
