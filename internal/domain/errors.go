package domain

import "errors"

var (
	ErrInternalServer        = errors.New("internal server error")
	ErrEmailDataValidation   = errors.New("email data validation error")
	ErrDublicateKeyViolation = errors.New("dublicate key violation error")
	ErrDataNotExists         = errors.New("data not exists error")
	ErrWrongAuthData         = errors.New("wrong auth data")
	ErrClientDataIncorrect   = errors.New("client data incorrect")
	ErrWrongLoginPassword    = errors.New("wrong login or password")
)
