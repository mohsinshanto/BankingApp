package account

import (
	appError "banking/errors"
	"banking/utils/response"

	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	service AccountService
}

func NewAccountController(service AccountService) *AccountController {
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
	idFloat, ok := userID.(float64)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "invalid user id type")
		return
	}
	id := uint(idFloat)
	result, err := ac.service.CreateAccount(id)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrAccountAlreadyExists):
			response.Error(c, http.StatusBadRequest, err.Error())

		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	response.Success(c, http.StatusCreated, "account created successfully", result)
}
func (ac *AccountController) Deposit(c *gin.Context) {
	var inputDepo Deposit
	if err := c.ShouldBindJSON(&inputDepo); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := ac.service.Deposit(inputDepo.AccountNo, inputDepo.Amount)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrAccountNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		case errors.Is(err, appError.ErrInvalidStatus),
			errors.Is(err, appError.ErrInvalidAmount):
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	response.Success(c, http.StatusOK, "money deposited successfully", result)
}
func (ac *AccountController) Withdraw(c *gin.Context) {
	var inputWithdraw Withdraw
	if err := c.ShouldBindJSON(&inputWithdraw); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := ac.service.Withdraw(inputWithdraw.AccountNo, inputWithdraw.Amount)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrAccountNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		case errors.Is(err, appError.ErrInvalidStatus),
			errors.Is(err, appError.ErrInvalidAmount),
			errors.Is(err, appError.ErrInsufficientBalance):
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	response.Success(c, http.StatusOK, "money withdrawn successfully", result)

}
func (ac *AccountController) MoneyTransfer(c *gin.Context) {
	var transferInput TransferInput
	if err := c.ShouldBindJSON(&transferInput); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := ac.service.MoneyTransfer(&transferInput)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrAccountNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		case errors.Is(err, appError.ErrSameAccountTransfer),
			errors.Is(err, appError.ErrSenderAccountNotActive),
			errors.Is(err, appError.ErrReceiverAccountNotActive),
			errors.Is(err, appError.ErrInsufficientBalance):
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
	var status AccountStatusUpdate
	if err := c.ShouldBindJSON(&status); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := ac.service.AccountStatusUpdate(accountNo, status.Status)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrAccountNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		case errors.Is(err, appError.ErrStatusAlreadySet),
			errors.Is(err, appError.ErrInvalidStatus):
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	response.Success(c, http.StatusOK, "Account status Updated", result)
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

	filter := TransactionFilter{
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
