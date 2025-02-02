package common_middlewares_errors

import "errors"

var (
	RoleNotFoundInContext = errors.New("role not found in Context")
	CurrentUserNotAdmin   = errors.New("current user not admin")
	CurrentUserNotClient  = errors.New("current user not client")
)
