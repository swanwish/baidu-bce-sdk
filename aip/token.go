package aip

import "time"

type Token struct {
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in"`
	Scope            string `json:"scope"`
	SessionKey       string `json:"session_key"`
	AccessToken      string `json:"access_token"`
	SessionSecret    string `json:"session_secret"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	CreateTime       int64  `json:"create_time"`
}

func (token Token) Valid() bool {
	if token.Error == "" {
		return time.Now().Unix()-token.CreateTime < token.ExpiresIn-TOKEN_PREFETCH_SECONDS
	}
	return false
}
