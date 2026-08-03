package main

import (
	"banking/config"
	"banking/internal/account"
	"banking/internal/user"
	"banking/models"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	if err := config.ConnectDB(); err != nil {
		log.Fatal(err)
	}
	if err := config.DB.AutoMigrate(&models.User{}, &models.Account{}, &models.Transaction{}); err != nil {
		log.Fatal(err)
	}
	// account dependency module injection
	accountModule := account.NewModule(config.DB)
	// user dependency module injection
	userModule := user.NewModule(config.DB)
	router := gin.Default()
	userModule.RegisterRoutes(router)
	accountModule.RegisterRoutes(router)
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}

}
