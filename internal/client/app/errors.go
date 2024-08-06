package app

import "errors"

var (
	ErrUserData             = errors.New("user data incorrect")
	ErrSrvInterationTimeout = errors.New("server not available")
)
