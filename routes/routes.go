package routes

import (
	"banking/controllers"
	"banking/middleware"

	"github.com/gin-gonic/gin"
)

func RouteHandler(r *gin.Engine, accountController *controllers.AccountController) {
	user := r.Group("/user")
	{
		user.POST("/register", controllers.Register)
		user.POST("/login", controllers.Login)
	}
	account := r.Group("/account")
	account.Use(middleware.AuthMiddleware)
	{
		account.POST("/", accountController.CreateAccount)
		account.POST("/deposit", accountController.Deposit)
		account.POST("/withdraw", accountController.Withdraw)
		account.POST("/transfer", accountController.MoneyTransfer)
		account.GET("/:accountNo/details", controllers.AccountDetails)
		account.GET("/transaction/:accountNo", controllers.GetTransactionsByAccount)
		account.GET("/:accountNo/summary", controllers.GetAccountSummary)
		account.GET("/:accountNo/stat", controllers.GetTransactionStatistics)
		account.PATCH("/:accountNo/status", controllers.AccountStatusUpdate)
	}

}
