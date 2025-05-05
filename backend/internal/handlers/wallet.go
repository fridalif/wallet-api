package handlers

import (
	"backend/internal/services"
	"backend/pkg/customerror"
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
	return walletHandler{
		WalletService: walletService,
	}
}

func (walletHandler walletHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/wallet", walletHandler.UpdateBalance)
	router.GET("/wallets/:id", walletHandler.UpdateBalance)
}
func (walletHandler walletHandler) GetBalance(ctx *gin.Context) {
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

func (walletHandler walletHandler) UpdateBalance(ctx *gin.Context) {

}
