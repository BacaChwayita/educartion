package cart

import "errors"

var (
	ErrInvalidQuantity   = errors.New("invalid quantity")
	ErrProductNotFound   = errors.New("product not found")
	ErrCartItemNotFound  = errors.New("cart item not found")
	ErrSessionTokenEmpty = errors.New("session token is required")
)
