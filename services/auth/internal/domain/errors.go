package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	UserNotFound          = errors.New("user not found")
	UserAlreadyExists     = errors.New("user already exists")
)
