package middleware

import (
	"banking/controllers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *gin.Context){
	authHeader:=c.GetHeader("Authorization")
	if authHeader == ""{
		c.JSON(http.StatusUnauthorized,gin.H{"error":"Authorization header is required"})
		c.Abort()
		return	
	}
	if !strings.HasPrefix(authHeader,"Bearer "){
		c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid authorization format"})
		c.Abort()
		return
	}
	tokenString := strings.TrimPrefix(authHeader,"Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return controllers.JwtSecret, nil
	})
	if err != nil{
		c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid token"})
		c.Abort()
		return
	}
	if !token.Valid{
		c.JSON(http.StatusUnauthorized,gin.H{"error":"Token is not valid"})
		c.Abort()
		return
	}
	// value,ok:=x.(int)
	claims,ok:= token.Claims.(jwt.MapClaims) 
	if !ok{
		c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid claims"})
		c.Abort()
		return
	} 
	c.Set("user-id",claims["user_id"])
	c.Next()
	
	
	
}
