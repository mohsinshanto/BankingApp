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

	account.GET("/details/:accountNo", m.Controller.AccountDetails)
	account.PUT("/status/:accountNo", m.Controller.AccountStatusUpdate)

	account.GET("/statistics/:accountNo", m.Controller.GetTransactionStatistics)
	account.GET("/summary/:accountNo", m.Controller.GetAccountSummary)
	account.GET("/transactions/:accountNo", m.Controller.GetTransactionsByAccount)
}
