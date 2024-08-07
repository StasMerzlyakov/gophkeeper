package domain

import "errors"

var (
	ErrServerInternal        = errors.New("internal server error")
	ErrNotAuthorized         = errors.New("user is not authorized")
	ErrDublicateKeyViolation = errors.New("dublicate key violation error")
	ErrDataNotExists         = errors.New("data not exists error")
	ErrAuthDataIncorrect     = errors.New("wrong auth data")

	ErrClientDataIncorrect = errors.New("client data incorrect")
	ErrClientInternal      = errors.New("client internal error")

	ErrClientDataIsNotRestored = errors.New("data is not restored")
	ErrServerIsNotResponding   = errors.New("server is not responding")

	ErrClientInteruptoin = errors.New("user interuption")
	ErrClientAppStopped  = errors.New("application stopped")
)
