package handler

type CtxKey int

const (
	RequestIDKey CtxKey = iota
	AuthorisationKey
)
