package dto

import "time"

type AccountDetails struct {
	AccountNo string    `json:"account_no"`
	Balance   float64   `json:"account_balance"`
	Status    string    `json:"account_status"`
	CreatedAt time.Time `json:"created_at"`
	UserName  string    `json:"username"`
	UserEmail string    `json:"email"`
}
