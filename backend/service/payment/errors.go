package payment

import "errors"

var (
	ErrInvalidPaymentRequest = errors.New("invalid payment request")
	ErrOrderNotFound         = errors.New("order not found")
	ErrPaymentNotFound       = errors.New("payment not found")
	ErrPaymentAlreadyExists  = errors.New("payment already exists for order")
)
