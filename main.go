package main

import (
	"banking/controllers"
	"banking/database"
	"banking/models"
	"banking/repositories"
	"banking/routes"
	"banking/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
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
	router.Run()

}
