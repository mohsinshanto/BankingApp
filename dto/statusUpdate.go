package dto

type AccountStatusUpdate struct {
	Status string `json:"status" binding:"required"`
}
type AccountStatusUpdateResponse struct {
	AccountNo string `json:"account_no"`
	Status    string `json:"status"`
}
