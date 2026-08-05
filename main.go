package main

// @title Banking API
// @version 1.0
// @description A Banking REST API built with Gin and GORM.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
import (
	"banking/config"
	"banking/internal/account"
	"banking/internal/user"
	"banking/models"
	"log"
	"os"

	_ "banking/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
