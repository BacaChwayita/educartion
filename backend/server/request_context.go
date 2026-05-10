package server

type CtxKey int

const (
	RequestIDKey CtxKey = iota
	AuthorisationKey
)
