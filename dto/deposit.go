package dto

type Deposit struct {
	AccountNo string  `json:"account_no" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
}
type DepositResponse struct {
	AccountNo string  `json:"account_no"`
	Balance   float64 `json:"balance"`
	Amount    float64 `json:"deposited_amount"`
	Status    string  `json:"status"`
}
