package response

import (
	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func Success(c *gin.Context, status int, message string, data any) {
	c.JSON(status, ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, ApiResponse{
		Success: false,
		Error:   message,
	})
}
