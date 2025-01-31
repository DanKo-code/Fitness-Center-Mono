package dtos

type RefreshResponse struct {
	AccessToken           string `json:"accessToken"`
	AccessTokenExpiration int    `json:"accessTokenExpiration"`
	User                  User   `json:"user"`
}
