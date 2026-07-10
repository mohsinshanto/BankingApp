package dto

type TransferInput struct{
	SenderAccNo string `json:"sender_acc" binding:"required"`
	ReceiverAccNo string `json:"receiver_acc" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}