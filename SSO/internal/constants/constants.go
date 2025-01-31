package constants

import "time"

type userRole struct {
	Client string
	Admin  string
	Coach  string
}

var ROLES = userRole{
	Client: "client",
	Admin:  "admin",
	Coach:  "coach",
}

type refreshToken struct {
	HttpOnly bool
	MaxAge   int
}

type cookieSettings struct {
	RefreshToken refreshToken
}

var COOKIE_SETTINGS = cookieSettings{
	RefreshToken: refreshToken{
		HttpOnly: true,
		MaxAge:   int(RefreshTokenExpiration),
	},
}

var AccessTokenExpiration = time.Minute * 15
var RefreshTokenExpiration = time.Hour * 24 * 7
