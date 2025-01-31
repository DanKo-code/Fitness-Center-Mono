package errors

import "errors"

var (
	VoidServiceData              = errors.New("void service data")
	ServiceAlreadyExists         = errors.New("service already exists")
	ServiceNotFound              = errors.New("service not found")
	CoachNotFound                = errors.New("coach not found")
	ReviewNotFound               = errors.New("coach not found")
	UserNotFound                 = errors.New("user not found")
	AbonementNotFound            = errors.New("abonement not found")
	InternalCoachServerError     = errors.New("internal coach server error")
	InternalUserServerError      = errors.New("internal user server error")
	InternalAbonementServerError = errors.New("internal abonement server error")
)
