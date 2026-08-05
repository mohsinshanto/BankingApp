package account

import (
	"banking/models"
	"time"
)

// withdraw
type Withdraw struct {
	AccountNo string  `json:"account_no" binding:"required" example:"ACC34218765"`
	Amount    float64 `json:"amount" binding:"required,gt=0" example:"500.00"`
}
type WithdrawResponse struct {
	AccountNo string  `json:"account_no"`
	Balance   float64 `json:"balance"`
	Amount    float64 `json:"withdrawn_amount"`
	Status    string  `json:"status"`
}

// transfer

type TransferInput struct {
	SenderAccNo   string  `json:"sender_acc" binding:"required" example:"ACC34218765"`
	ReceiverAccNo string  `json:"receiver_acc" binding:"required" example:"ACC34211234"`
	Amount        float64 `json:"amount" binding:"required,gt=0" example:"500.00"`
}
type TransferResponse struct {
	SenderAccountNo   string  `json:"sender_account_no"`
	ReceiverAccountNo string  `json:"receiver_account_no"`
	TransferAmount    float64 `json:"transfer_amount"`
	SenderNewBalance  float64 `json:"sender_new_balance"`
}

// TransactionStat

type TransactionStatisticsResponse struct {
	AccountNo            string `json:"account_no"`
	TodayTransaction     int64  `json:"today_transaction"`
	ThisWeekTransaction  int64  `json:"this_week_transaction"`
	ThisMonthTransaction int64  `json:"this_month_transaction"`
	TotalTransaction     int64  `json:"total_transaction"`
}

// Transactionfilter

type TransactionFilter struct {
	Page            string
	Limit           string
	TransactionType string
	FromDate        string
	ToDate          string
	SortBy          string
	From            time.Time
	To              time.Time
	OrderBy         string
	PageNo          int
	PageSize        int
	Offset          int
}
type TransactionListResponse struct {
	AccountNo         string               `json:"account_no"`
	CurrentPage       int                  `json:"current_page"`
	Limit             int                  `json:"limit"`
	TotalPages        int64                `json:"total_pages"`
	TotalTransactions int64                `json:"total_transactions"`
	Transactions      []models.Transaction `json:"transactions"`
}
type TransactionQueryResult struct {
	Total        int64
	Transactions []models.Transaction
}

// summary

type AccountSummaryResponse struct {
	AccountNo             string  `json:"account_no"`
	CurrentBalance        float64 `json:"current_balance"`
	TotalTransactions     int64   `json:"total_transactions"`
	TotalDeposit          float64 `json:"total_deposit"`
	TotalWithdraw         float64 `json:"total_withdraw"`
	TotalTransferSent     float64 `json:"total_transfer_sent"`
	TotalTransferReceived float64 `json:"total_transfer_received"`
}
type AccountSummary struct {
	TotalTransactions     int64   `json:"total_transactions"`
	TotalDeposit          float64 `json:"total_deposit"`
	TotalWithdraw         float64 `json:"total_withdraw"`
	TotalTransferSent     float64 `json:"total_transfer_sent"`
	TotalTransferReceived float64 `json:"total_transfer_received"`
}

// statusUpdate

type AccountStatusUpdate struct {
	Status string `json:"status" binding:"required" example:"BLOCKED"`
}
type AccountStatusUpdateResponse struct {
	AccountNo string `json:"account_no"`
	Status    string `json:"status"`
}

// deposit

type Deposit struct {
	AccountNo string  `json:"account_no" binding:"required" example:"ACC87112088"`
	Amount    float64 `json:"amount" binding:"required,gt=0" example:"500.00"`
}
type DepositResponse struct {
	AccountNo string  `json:"account_no"`
	Balance   float64 `json:"balance"`
	Amount    float64 `json:"deposited_amount"`
	Status    string  `json:"status"`
}

// accountDetails

type AccountDetailsResponse struct {
	AccountNo string    `json:"account_no"`
	Balance   float64   `json:"account_balance"`
	Status    string    `json:"account_status"`
	CreatedAt time.Time `json:"created_at"`
	UserName  string    `json:"username"`
	UserEmail string    `json:"email"`
}

// account

type AccountResponse struct {
	AccountNo string  `json:"account_no"`
	Balance   float64 `json:"balance"`
	Status    string  `json:"status"`
}
