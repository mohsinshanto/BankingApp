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

// CreateAccount godoc
// @Summary Create a new bank account
// @Description Creates a bank account for the authenticated user.
// @Tags Account
// @Security BearerAuth
// @Produce json
// @Success 201 {object} response.ApiResponse "Account created successfully"
// @Failure 400 {object} response.ApiResponse "Account already exists"
// @Failure 401 {object} response.ApiResponse "Unauthorized"
// @Failure 500 {object} response.ApiResponse "Internal server error"
// @Router /account/ [post]
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
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	response.Success(c, http.StatusCreated, "account created successfully", result)
}

// Deposit godoc
// @Summary Deposit money into an account
// @Description Deposit money into an active bank account.
// @Tags Account
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body account.Deposit true "Deposit Request"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Failure 401 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Router /account/deposit [post]
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
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	response.Success(c, http.StatusOK, "money deposited successfully", result)
}

// Withdraw godoc
// @Summary Withdraw money from an account
// @Description Withdraw money from an active bank account
// @Tags Account
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body account.Withdraw true "Withdraw Request"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Failure 401 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Failure 500 {object} response.ApiResponse
// @Router /account/withdraw [post]
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
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	response.Success(c, http.StatusOK, "money withdrawn successfully", result)

}

// MoneyTransfer godoc
// @Summary Transfer money between two accounts
// @Description Transfer money from one active account to another active account.
// @Tags Account
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body account.TransferInput true "Money Transfer Request"
// @Success 200 {object} response.ApiResponse "Money transferred successfully"
// @Failure 400 {object} response.ApiResponse "Invalid request, inactive account, insufficient balance, or same account transfer"
// @Failure 401 {object} response.ApiResponse "Unauthorized"
// @Failure 404 {object} response.ApiResponse "Account not found"
// @Failure 500 {object} response.ApiResponse "Internal server error"
// @Router /account/transfer [post]
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
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	response.Success(c, http.StatusOK, "Money transferred successfully", result)

}

// AccountDetails godoc
// @Summary Get account details
// @Description Retrieve details of a bank account by account number.
// @Tags Account
// @Security BearerAuth
// @Produce json
// @Param accountNo path string true "Account Number"
// @Success 200 {object} response.ApiResponse "Account details retrieved successfully"
// @Failure 401 {object} response.ApiResponse "Unauthorized"
// @Failure 404 {object} response.ApiResponse "Account not found"
// @Failure 500 {object} response.ApiResponse "Internal server error"
// @Router /account/details/{accountNo} [get]
func (ac *AccountController) AccountDetails(c *gin.Context) {
	accountNo := c.Param("accountNo")
	details, err := ac.service.AccountDetails(accountNo)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrAccountNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	response.Success(c, http.StatusOK, "Account Details Found", details)
}

// AccountStatusUpdate godoc
// @Summary Update account status
// @Description Update the status of an existing bank account.
// @Tags Account
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param accountNo path string true "Account Number"
// @Param request body account.AccountStatusUpdate true "Account Status Update Request"
// @Success 200 {object} response.ApiResponse "Account status updated successfully"
// @Failure 400 {object} response.ApiResponse "Invalid request, invalid status, or status already set"
// @Failure 401 {object} response.ApiResponse "Unauthorized"
// @Failure 404 {object} response.ApiResponse "Account not found"
// @Failure 500 {object} response.ApiResponse "Internal server error"
// @Router /account/status/{accountNo} [put]
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
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	response.Success(c, http.StatusOK, "Account status Updated", result)
}

// GetTransactionStatistics godoc
// @Summary Get transaction statistics
// @Description Retrieve transaction statistics for a specific bank account.
// @Tags Account
// @Security BearerAuth
// @Produce json
// @Param accountNo path string true "Account Number"
// @Success 200 {object} response.ApiResponse "Transaction statistics retrieved successfully"
// @Failure 401 {object} response.ApiResponse "Unauthorized"
// @Failure 404 {object} response.ApiResponse "Account not found"
// @Failure 500 {object} response.ApiResponse "Internal server error"
// @Router /account/statistics/{accountNo} [get]
func (ac *AccountController) GetTransactionStatistics(c *gin.Context) {
	accountNo := c.Param("accountNo")
	result, err := ac.service.GetTransactionStat(accountNo)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrAccountNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	response.Success(c, http.StatusOK, "transaction statistics retrieved successfully", result)
}

// GetAccountSummary godoc
// @Summary Get account summary
// @Description Retrieve a summary of a bank account, including balance and transaction totals.
// @Tags Account
// @Security BearerAuth
// @Produce json
// @Param accountNo path string true "Account Number"
// @Success 200 {object} response.ApiResponse "Account summary retrieved successfully"
// @Failure 401 {object} response.ApiResponse "Unauthorized"
// @Failure 404 {object} response.ApiResponse "Account not found"
// @Failure 500 {object} response.ApiResponse "Internal server error"
// @Router /account/summary/{accountNo} [get]
func (ac *AccountController) GetAccountSummary(c *gin.Context) {
	accountNo := c.Param("accountNo")
	result, err := ac.service.GetAccountSummary(accountNo)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrAccountNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	response.Success(c, http.StatusOK, "Account summary retrieved successfully", result)
}

// GetTransactionsByAccount godoc
// @Summary Get account transactions
// @Description Retrieve transactions for a specific account with pagination, filtering, date range, and sorting.
// @Tags Account
// @Security BearerAuth
// @Produce json
// @Param accountNo path string true "Account Number"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of records per page" default(5)
// @Param type query string false "Transaction type (DEPOSIT, WITHDRAW, TRANSFER)"
// @Param from query string false "Start date (YYYY-MM-DD), e.g. 2026-08-01"
// @Param to query string false "End date (YYYY-MM-DD), e.g. 2026-08-31"
// @Param sort query string false "Sort order (newest, oldest)" default(newest)
// @Success 200 {object} response.ApiResponse "Transactions retrieved successfully"
// @Failure 400 {object} response.ApiResponse "Invalid query parameters"
// @Failure 401 {object} response.ApiResponse "Unauthorized"
// @Failure 404 {object} response.ApiResponse "Account not found"
// @Failure 500 {object} response.ApiResponse "Internal server error"
// @Router /account/transactions/{accountNo} [get]
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
			response.Error(c, http.StatusInternalServerError, "internal server error")

		}
		return
	}
	response.Success(c, http.StatusOK, "transaction retrieved successfully", result)
}
