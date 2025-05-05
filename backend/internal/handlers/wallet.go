package handlers

import "github.com/gin-gonic/gin"

type WalletHandlerI interface {
	RegisterRoutes(router *gin.RouterGroup)
	GetBalance(ctx *gin.Context)
	UpdateBalance(ctx *gin.Context)
}
