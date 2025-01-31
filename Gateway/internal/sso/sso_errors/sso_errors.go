package sso_errors

import "errors"

var (
	FingerPrintNotFoundInContext = errors.New("fingerprint not found in Context")
	FingerprintIsNotValidString  = errors.New("fingerprint is not a valid string")
)
