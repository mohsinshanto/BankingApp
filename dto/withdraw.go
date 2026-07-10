package dto
type Withdraw struct{
	AccountNo string `json:"account_no" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
}