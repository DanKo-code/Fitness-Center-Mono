package dtos

import (
	userGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
)

type SignInResponseCommand struct {
	AccessToken            string `json:"access_token"`
	RefreshToken           string `json:"refresh_token"`
	AccessTokenExpiration  string `json:"access_token_expiration"`
	RefreshTokenExpiration string `json:"refresh_token_expiration"`
	User                   *userGRPC.UserObject
}
