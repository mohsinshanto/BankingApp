package dto

type TransferInput struct {
	SenderAccNo   string  `json:"sender_acc" binding:"required"`
	ReceiverAccNo string  `json:"receiver_acc" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
}
type TransferResponse struct {
	SenderAccountNo   string  `json:"sender_account_no"`
	ReceiverAccountNo string  `json:"receiver_account_no"`
	TransferAmount    float64 `json:"transfer_amount"`
	SenderNewBalance  float64 `json:"sender_new_balance"`
}
