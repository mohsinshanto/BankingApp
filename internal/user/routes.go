package user

import "github.com/gin-gonic/gin"

func (m *Module) RegisterRoutes(r *gin.Engine) {
	user := r.Group("/users")

	user.POST("/register", m.Controller.Register)
	user.POST("/login", m.Controller.Login)
}
