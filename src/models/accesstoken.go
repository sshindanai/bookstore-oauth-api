package models

import (
	"github.com/sshindanai/repo/bookstore-oauth-api/src/utils/errors"
)

const (
	expirationTime            = 24
	grantTypePassword         = "password"
	grantTypeClientCredential = "client_credentials"
)

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`

	// Used for password grant type
	Username string `json:"username"`
	Password string `json:"password"`

	// Used for client_credentials grant type
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
	ClientID    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

type AccessTokenConcurrent struct {
	Result *AccessToken
	Error  *errors.RestErr
}
