package auth

import "errors"

// These are the errors that exposed auth package functions can return

var (
	ErrHashPassword      = errors.New("could not hash password")
	ErrDBInsert          = errors.New("could not insert into database")
	ErrAccountNotFound   = errors.New("account not found")
	ErrAccountInactive   = errors.New("account inactive")
	ErrPasswordIncorrect = errors.New("password incorrect")
	ErrUnkown            = errors.New("unkown error")
	ErrFailedToCreateJWT = errors.New("JWT not created")
)
