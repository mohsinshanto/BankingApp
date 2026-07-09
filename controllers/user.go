package controllers

import (
	"banking/database"
	"banking/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)
var JwtSecret = []byte("mohsin5702shanto")
func Register(c *gin.Context){
	var user models.RegistrationInput
	if err := c.ShouldBindJSON(&user); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	hashpassword,err:=bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost)
	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	dbUser := models.User{
		Name: user.Name,
		Email: user.Email,
		Password: string(hashpassword),
	}
	if err:=database.DB.Create(&dbUser).Error; err !=nil{
        c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	c.JSON(http.StatusCreated,gin.H{"msg":"User has been created","user-id":dbUser.ID})
}
func Login(c *gin.Context){
	var loginInput models.LoginInput
	if err:= c.ShouldBindJSON(&loginInput); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	var user models.User
	if err:= database.DB.Where("email=?",loginInput.Email).First(&user).Error; err != nil{
        c.JSON(http.StatusUnauthorized,gin.H{"error":"invalid email or password"})
		return
	}
	if err:=bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(loginInput.Password)); err != nil{
		c.JSON(http.StatusUnauthorized, gin.H{ "error": "invalid email or password",})
        return
	}
	claims :=jwt.MapClaims{
		"user_id": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	tokenString,err:=token.SignedString(JwtSecret)
	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	c.JSON(http.StatusOK,gin.H{"your token":tokenString})

	}
