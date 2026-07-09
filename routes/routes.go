package routes

import (
	"banking/controllers"
	"banking/middleware"

	"github.com/gin-gonic/gin"
)


func RouteHandler(r *gin.Engine){
	user := r.Group("/user") 
	{
		user.POST("/register",controllers.Register)
		user.POST("/login",controllers.Login)
	}
	account:= r.Group("/account")
	account.Use(middleware.AuthMiddleware)
	{
     account.POST("/",controllers.CreateAccount)
	 account.POST("/deposit",controllers.Deposit)
	}
	
}