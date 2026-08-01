package main

import (
	"banking/config"
	"banking/controllers"
	"banking/database"
	"banking/models"
	"banking/repositories"
	"banking/routes"
	"banking/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	if err := database.ConnectDB(); err != nil {
		log.Fatal(err)
	}
	if err := database.DB.AutoMigrate(&models.User{}, &models.Account{}, &models.Transaction{}); err != nil {
		log.Fatal(err)
	}
	// account dependency injection
	accountRepo := repositories.NewAccountRepository(database.DB)
	accountService := services.NewAccountService(accountRepo)
	accountController := controllers.NewAccountController(accountService)
	// user dependency injection
	userRepo := repositories.NewUserRepository(database.DB)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)
	router := gin.Default()
	routes.RouteHandler(router, accountController, userController)
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}

}
