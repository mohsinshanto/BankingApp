package routes

import (
	"banking/controllers"
	"banking/middleware"

	"github.com/gin-gonic/gin"
)

func RouteHandler(r *gin.Engine, accountController *controllers.AccountController, userController *controllers.UserController) {
	user := r.Group("/user")
	{
		user.POST("/register", userController.Register)
		user.POST("/login", userController.Login)
	}
	account := r.Group("/account")
	account.Use(middleware.AuthMiddleware)
	{
		account.POST("/", accountController.CreateAccount)
		account.POST("/deposit", accountController.Deposit)
		account.POST("/withdraw", accountController.Withdraw)
		account.POST("/transfer", accountController.MoneyTransfer)
		account.GET("/:accountNo/details", accountController.AccountDetails)
		account.GET("/transaction/:accountNo", controllers.GetTransactionsByAccount)
		account.GET("/:accountNo/summary", controllers.GetAccountSummary)
		account.GET("/:accountNo/stat", accountController.GetTransactionStatistics)
		account.PATCH("/:accountNo/status", accountController.AccountStatusUpdate)
	}

}
