package dto
type AccountSummaryResponse struct{
	AccountNo  string  `json:"account_no"`
	CurrentBalance float64 `json:"current_balance"`
	TotalTransactions int64 `json:"total_transactions"`
	TotalDeposit float64  `json:"total_deposit"`
	TotalWtihdraw float64 `json:"total_withdraw"`
	TotalTransferSent float64 `json:"total_transfer_sent"`
	TotalTransferReceived float64 `json:"total_transfer_received"`

}