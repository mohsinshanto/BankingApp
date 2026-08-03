package account

import (
	"banking/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (m *Module) RegisterRoutes(r *gin.Engine) {
	account := r.Group("/account")
	account.Use(middleware.AuthMiddleware)

	account.POST("/", m.Controller.CreateAccount)
	account.POST("/deposit", m.Controller.Deposit)
	account.POST("/withdraw", m.Controller.Withdraw)
	account.POST("/transfer", m.Controller.MoneyTransfer)

	account.GET("/:accountNo", m.Controller.AccountDetails)
	account.PUT("/:accountNo/status", m.Controller.AccountStatusUpdate)

	account.GET("/:accountNo/statistics", m.Controller.GetTransactionStatistics)
	account.GET("/:accountNo/summary", m.Controller.GetAccountSummary)
	account.GET("/:accountNo/transactions", m.Controller.GetTransactionsByAccount)
}
