package dtos

import userGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"

type SignUpResponseCommand struct {
	AccessToken            string
	RefreshToken           string
	AccessTokenExpiration  string
	RefreshTokenExpiration string
	User                   *userGRPC.UserObject
}
