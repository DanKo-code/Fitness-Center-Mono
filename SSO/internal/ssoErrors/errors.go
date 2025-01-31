package ssoErrors

import "errors"

var (
	UserAlreadyExists   = errors.New("user already exists")
	InvalidPassword     = errors.New("invalid password")
	InvalidAccessToken  = errors.New("invalid access token")
	InvalidRefreshToken = errors.New("invalid refresh token")
	InvalidFingerPrint  = errors.New("invalid fingerprint")

	UserNotFound           = errors.New("user not found")
	RefreshSessionNotFound = errors.New("refresh session not found")

	VoidUserData = errors.New("void userData")
)
