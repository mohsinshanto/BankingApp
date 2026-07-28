package error

import "errors"

var (
	ErrAccountAlreadyExists     = errors.New("account already exists")
	ErrAccountNotFound          = errors.New("account not found")
	ErrInvalidStatus            = errors.New("invalid account status")
	ErrInsufficientBalance      = errors.New("Insufficient balance")
	ErrInvalidAmount            = errors.New("invalid amount")
	ErrSameAccountTransfer      = errors.New("cannot transfer to the same account")
	ErrSenderAccountNotActive   = errors.New("sender account is not active")
	ErrReceiverAccountNotActive = errors.New("receiver account is not active")
)
