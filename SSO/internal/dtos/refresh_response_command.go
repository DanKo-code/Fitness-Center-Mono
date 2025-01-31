package dtos

import ssoGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"

type RefreshResponseCommand struct {
	AccessToken            string
	RefreshToken           string
	AccessTokenExpiration  string
	RefreshTokenExpiration string
	User                   *ssoGRPC.UserObject
}
