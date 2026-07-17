package dto

type AccountStatusUpdate struct {
	Status string `json:"status" binding:"required"`
}
