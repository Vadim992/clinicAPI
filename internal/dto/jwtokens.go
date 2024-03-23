package dto

import "encoding/json"

type JWTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewJWTokens(accessToken, refreshToken string) *JWTokens {
	return &JWTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func (jwt *JWTokens) EncodeToJSON() ([]byte, error) {
	b, err := json.Marshal(jwt)

	if err != nil {
		return nil, err
	}

	return b, nil
}
