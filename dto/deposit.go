package dto

type Deposit struct{
	AccountNo string `json:"account_no" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
}