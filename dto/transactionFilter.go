package dto

import (
	"banking/models"
	"time"
)

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
