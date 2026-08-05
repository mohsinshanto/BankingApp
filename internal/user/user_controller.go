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

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags User
// @Accept json
// @Produce json
// @Param request body user.RegistrationInput true "User Registration"
// @Success 201 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Router /user/register [post]
func (uc *UserController) Register(c *gin.Context) {
	var registrationInput RegistrationInput
	if err := c.ShouldBindJSON(&registrationInput); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := uc.service.Register(&registrationInput)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrEmailAlreadyExists):
			response.Error(c, http.StatusBadRequest, err.Error())

		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	response.Success(c, http.StatusCreated, "new user created", result)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password, then generate a JWT token.
// @Tags User
// @Accept json
// @Produce json
// @Param request body user.LoginInput true "Login credentials"
// @Success 200 {object} response.ApiResponse "Token generated successfully"
// @Failure 400 {object} response.ApiResponse "Invalid request body"
// @Failure 401 {object} response.ApiResponse "Invalid email or password"
// @Failure 500 {object} response.ApiResponse "Internal server error"
// @Router /user/login [post]
func (uc *UserController) Login(c *gin.Context) {
	var loginInput LoginInput
	if err := c.ShouldBindJSON(&loginInput); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := uc.service.Login(&loginInput)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrInvalidCredentials):
			response.Error(c, http.StatusUnauthorized, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	response.Success(c, http.StatusOK, "token generated successfully", result)
}
