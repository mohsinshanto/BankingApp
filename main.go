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
	repo := repositories.NewAccountRepository(database.DB)
	service := services.NewAccountService(repo)
	accountController := controllers.NewAccountController(service)
	router := gin.Default()
	routes.RouteHandler(router, accountController)
	router.Run()

}
