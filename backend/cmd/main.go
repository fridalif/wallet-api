package main

import (
	"backend/internal/handlers"
	"backend/internal/repos"
	"backend/internal/services"
	"backend/pkg/config"
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.NewConfig("config.env")
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	walletRepository, err := repos.NewWalletRepository(config)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	defer walletRepository.ClosePull()
	err = walletRepository.CreateTables(context.Background())
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	walletService := services.NewWalletService(walletRepository)
	walletHandlers := handlers.NewWalletHandler(walletService)

	router := gin.Default()
	api := router.Group("/api")
	v1 := api.Group("/v1")
	walletHandlers.RegisterRoutes(v1)

	router.Run(fmt.Sprintf("%s:%s", config.WebHost, config.WebPort))
}
