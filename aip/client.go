package aip

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/swanwish/go-common/logs"
	"github.com/swanwish/go-common/utils"
)

const (
	AIP_TOKEN_URL           = "https://aip.baidubce.com/oauth/2.0/token"
	GRANT_TYPE              = "client_credentials"
	PARAM_KEY_GRANT_TYPE    = "grant_type"
	PARAM_KEY_CLIENT_ID     = "client_id"
	PARAM_KEY_CLIENT_SECRET = "client_secret"
)

type TokenDao interface {
	GetToken(clientId, clientSecret string) (Token, error)
	SaveToken(clientId, clientSecret string, token Token) error
}

type Client struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Token        *Token `json:"token"`
	TokenDao     TokenDao
}

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

func NewClient(clientId, clientSecret string) *Client {
	return &Client{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
}

func (client *Client) GetAccessToken() (string, error) {
	if client.ClientId == "" || client.ClientSecret == "" {
		logs.Errorf("The client id or client secret not specified")
		return "", ErrInvalidParameter
	}
	if client.Token != nil && client.Token.Valid() {
		return client.Token.AccessToken, nil
	}
	if client.TokenDao != nil {
		token, err := client.TokenDao.GetToken(client.ClientId, client.ClientSecret)
		if err != nil {
			logs.Errorf("Failed to get token from token dao, the error is %#v", err)
		} else {
			client.Token = &token
			return token.AccessToken, nil
		}
	}
	data := url.Values{}
	data.Set(PARAM_KEY_GRANT_TYPE, GRANT_TYPE)
	data.Set(PARAM_KEY_CLIENT_ID, client.ClientId)
	data.Set(PARAM_KEY_CLIENT_SECRET, client.ClientSecret)
	status, content, err := utils.PostRequest(AIP_TOKEN_URL, data, nil)
	if err != nil {
		logs.Errorf("Failed to get token, the error is %#v", err)
		return "", err
	}
	if status != http.StatusOK {
		logs.Errorf("Failed to get token, the status code is %d", status)
		return "", ErrInvalidStatus
	}
	token := Token{}
	err = json.Unmarshal(content, &token)
	if err != nil {
		logs.Errorf("Failed to unmarshal token, the error is %#v", err)
		return "", err
	}
	token.CreateTime = time.Now().Unix()
	client.Token = &token
	if client.TokenDao != nil {
		client.TokenDao.SaveToken(client.ClientId, client.ClientSecret, token)
	}
	if token.Valid() {
		return token.AccessToken, err
	}
	return "", ErrInvalidToken
}

func (token Token) Valid() bool {
	if token.Error == "" {
		return time.Now().Unix()-token.CreateTime > token.ExpiresIn
	}
	return false
}
