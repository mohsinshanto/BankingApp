package user

import (
	appError "banking/errors"
	"banking/utils/response"

	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service UserService
}

func NewUserController(service UserService) *UserController {
	return &UserController{
		service: service,
	}
}
func (uc *UserController) Register(c *gin.Context) {
	var registrationInput RegistrationInput
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
	var loginInput LoginInput
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
