package handlers

import (
	"backend/internal/services"
	"backend/pkg/customerror"
	"backend/pkg/requests"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type WalletHandlerI interface {
	RegisterRoutes(router *gin.RouterGroup)
	GetBalance(ctx *gin.Context)
	UpdateBalance(ctx *gin.Context)
}

type walletHandler struct {
	WalletService services.WalletServiceI
}

func NewWalletHandler(walletService services.WalletServiceI) WalletHandlerI {
	return &walletHandler{
		WalletService: walletService,
	}
}

func (walletHandler *walletHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/wallet", walletHandler.UpdateBalance)
	router.GET("/wallets/:id", walletHandler.UpdateBalance)
}
func (walletHandler *walletHandler) GetBalance(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"data":   gin.H{},
			"error":  "Wrong uuid",
		})
		return
	}
	balance, err := walletHandler.WalletService.GetBalance(id)
	if err == pgx.ErrNoRows {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusNotFound,
			"data":   gin.H{},
			"error":  "Wallet not found",
		})
		return
	}
	if err != nil {
		customError := err.(customerror.CustomError)
		customError.AppendModule("GetBalance")
		log.Printf("%s", customError.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"data":   gin.H{},
			"error":  "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": gin.H{
			"balance": balance,
		},
		"error": nil,
	})
}

func (walletHandler *walletHandler) UpdateBalance(ctx *gin.Context) {
	var userRequest requests.UpdateBalanceRequest
	err := ctx.ShouldBindJSON(&userRequest)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"body":   gin.H{},
			"error":  "Wrong input",
		})
		return
	}
	err = walletHandler.WalletService.UpdateBalance(userRequest.WalletId, userRequest.OperationType, userRequest.Amount)
	if err == customerror.ErrWrongAmount {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"body":   gin.H{},
			"error":  "Amount cant be less than zero",
		})
		return
	}
	if err == customerror.ErrWrongOperation {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"body":   gin.H{},
			"error":  "Operation must be DEPOSIT or WITHDRAW",
		})
		return
	}
	if err == pgx.ErrNoRows {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusNotFound,
			"body":   gin.H{},
			"error":  "Wallet not found",
		})
		return
	}
	if err != nil {
		customError := err.(customerror.CustomError)
		customError.AppendModule("UpdateBalance")
		log.Printf("%s", customError.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status": http.StatusInternalServerError,
			"data":   gin.H{},
			"error":  "Internal Server Error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"body":   gin.H{},
		"error":  nil,
	})
}
