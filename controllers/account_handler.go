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

type AccountController struct {
	service services.AccountService
}

func NewAccountController(service services.AccountService) *AccountController {
	return &AccountController{
		service: service,
	}
}
func (ac *AccountController) CreateAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "invalid user id type")
		return
	}
	account, err := ac.service.CreateAccount(id)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrAccountAlreadyExists):
			response.Error(c, http.StatusBadRequest, err.Error())

		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	response.Success(c, http.StatusCreated, "account created successfully", account)
}
func (ac *AccountController) Deposit(c *gin.Context) {
	var inputDepo dto.Deposit
	if err := c.ShouldBindJSON(&inputDepo); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	account, err := ac.service.Deposit(inputDepo.AccountNo, inputDepo.Amount)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrAccountNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		case errors.Is(err, appError.ErrInvalidStatus):
			response.Error(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, appError.ErrInvalidAmount):
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	response.Success(c, http.StatusOK, "money deposited successfully", account)
}
func (ac *AccountController) Withdraw(c *gin.Context) {
	var inputWithdraw dto.Withdraw
	if err := c.ShouldBindJSON(&inputWithdraw); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	account, err := ac.service.Withdraw(inputWithdraw.AccountNo, inputWithdraw.Amount)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrAccountNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		case errors.Is(err, appError.ErrInvalidStatus):
			response.Error(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, appError.ErrInvalidAmount):
			response.Error(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, appError.ErrInsufficientBalance):
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	response.Success(c, http.StatusOK, "money withdrawn successfully", account)

}
func (ac *AccountController) MoneyTransfer(c *gin.Context) {
	var transferInput dto.TransferInput
	if err := c.ShouldBindJSON(&transferInput); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := ac.service.MoneyTransfer(&transferInput)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrAccountNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		case errors.Is(err, appError.ErrSameAccountTransfer):
			response.Error(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, appError.ErrSenderAccountNotActive):
			response.Error(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, appError.ErrReceiverAccountNotActive):
			response.Error(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, appError.ErrInsufficientBalance):
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	response.Success(c, http.StatusOK, "Money transffered successfully", result)

}
