package dto

type AccountResponse struct {
	AccountNo string  `json:"account_no"`
	Balance   float64 `json:"balance"`
	Status    string  `json:"status"`
}
