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
	userID, exists := c.Get("user-id")
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
func (ac *AccountController) AccountDetails(c *gin.Context) {
	accountNo := c.Param("accountNo")
	details, err := ac.service.AccountDetails(accountNo)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrAccountNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	response.Success(c, http.StatusOK, "Account Details Found", details)
}
func (ac *AccountController) AccountStatusUpdate(c *gin.Context) {
	accountNo := c.Param("accountNo")
	var status dto.AccountStatusUpdate
	if err := c.ShouldBindJSON(&status); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	account, err := ac.service.AccountStatusUpdate(accountNo, status.Status)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrAccountNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		case errors.Is(err, appError.ErrStatusAlreadySet):
			response.Error(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, appError.ErrInvalidStatus):
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	response.Success(c, http.StatusOK, "Account status Updated", account)
}
func (ac *AccountController) GetTransactionStatistics(c *gin.Context) {
	accountNo := c.Param("accountNo")
	result, err := ac.service.GetTransactionStat(accountNo)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrAccountNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	response.Success(c, http.StatusOK, "transaction retrieved successfully", result)
}
func (ac *AccountController) GetAccountSummary(c *gin.Context) {
	accountNo := c.Param("accountNo")
	result, err := ac.service.GetAccountSummary(accountNo)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrAccountNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	response.Success(c, http.StatusOK, "here is the summary", result)
}
func (ac *AccountController) GetTransactionsByAccount(c *gin.Context) {
	accountNo := c.Param("accountNo")

	filter := dto.TransactionFilter{
		Page:            c.DefaultQuery("page", "1"),
		Limit:           c.DefaultQuery("limit", "5"),
		TransactionType: c.Query("type"),
		FromDate:        c.Query("from"),
		ToDate:          c.Query("to"),
		SortBy:          c.DefaultQuery("sort", "newest"),
	}

	result, err := ac.service.GetTransactionsByAccount(accountNo, &filter)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrAccountNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		case errors.Is(err, appError.ErrInvalidTransactionType),
			errors.Is(err, appError.ErrInvalidDate),
			errors.Is(err, appError.ErrInvalidDateRange),
			errors.Is(err, appError.ErrInvalidSortOption):
			response.Error(c, http.StatusBadRequest, err.Error())

		default:
			response.Error(c, http.StatusInternalServerError, err.Error())

		}
		return
	}
	response.Success(c, http.StatusOK, "transaction retrieved successfully", result)
}
