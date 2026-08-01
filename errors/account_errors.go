package error

import "errors"

var (
	ErrAccountAlreadyExists     = errors.New("account already exists")
	ErrAccountNotFound          = errors.New("account not found")
	ErrInvalidStatus            = errors.New("invalid account status")
	ErrStatusAlreadySet         = errors.New("account is already in this status")
	ErrInsufficientBalance      = errors.New("Insufficient balance")
	ErrInvalidAmount            = errors.New("invalid amount")
	ErrSameAccountTransfer      = errors.New("cannot transfer to the same account")
	ErrSenderAccountNotActive   = errors.New("sender account is not active")
	ErrReceiverAccountNotActive = errors.New("receiver account is not active")
	ErrInvalidTransactionType   = errors.New("invalid transaction type")
	ErrInvalidDate              = errors.New("invalid date (YYYY-MM-DD) expected")
	ErrInvalidDateRange         = errors.New("invalid date range")
	ErrInvalidSortOption        = errors.New("invalid sorting option")
)
