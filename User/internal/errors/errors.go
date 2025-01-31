package errors

import "errors"

var (
	UserNotFound    = errors.New("user not found")
	InvalidPassword = errors.New("invalid password")
)
