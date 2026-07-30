package controllers

import (
	"banking/dto"
	appError "banking/errors"
	"banking/response"
	"banking/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{
		service: service,
	}
}
func (uc *UserController) Register(c *gin.Context) {
	var registrationInput dto.RegistrationInput
	if err := c.ShouldBindJSON(&registrationInput); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	err := uc.service.Register(&registrationInput)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "new user created", " data is saved")
}
func (uc *UserController) Login(c *gin.Context) {
	var loginInput dto.LoginInput
	if err := c.ShouldBindJSON(&loginInput); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := uc.service.Login(&loginInput)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrInvalidCredentials):
			response.Error(c, http.StatusUnauthorized, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
	}
	response.Success(c, http.StatusOK, "token generated successfully", token)
}
