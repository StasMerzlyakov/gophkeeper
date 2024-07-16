package domain

import "errors"

var (
	ErrEmailDataValidation   = errors.New("email data validation error")
	ErrDublicateKeyViolation = errors.New("dublicate key violation error")
	ErrDataNotExists         = errors.New("data not exists error")
)
