package cases

import "errors"

var (
	ErrInternal          = errors.New("Server internal error")
	ErrLoginFailed       = errors.New("Login failed")
	ErrEmailAlreadyTaken = errors.New("Email already taken")
)
