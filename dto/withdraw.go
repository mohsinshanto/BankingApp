package dto

type Withdraw struct {
	AccountNo string  `json:"account_no" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
}
type WithdrawResponse struct {
	AccountNo string  `json:"account_no"`
	Balance   float64 `json:"balance"`
	Amount    float64 `json:"withdrawn_amount"`
	Status    string  `json:"status"`
}
