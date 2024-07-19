package domain

import "errors"

var (
	ErrServerInternal        = errors.New("internal server error")
	ErrEmailDataValidation   = errors.New("email data validation error")
	ErrDublicateKeyViolation = errors.New("dublicate key violation error")
	ErrDataNotExists         = errors.New("data not exists error")
	ErrAuthDataIncorrect     = errors.New("wrong auth data")
	ErrClientDataIncorrect   = errors.New("client data incorrect")
)
