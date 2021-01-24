package models

import "github.com/sshindanai/bookstore-utils-go/resterrors"

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`

	// Used for password grant type
	Email    string `json:"email"`
	Password string `json:"password"`

	// Used for client_credentials grant type
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	UserID      int64  `json:"user_id"`
	ClientID    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

type AccessTokenConcurrent struct {
	Result *AccessToken
	Error  *resterrors.RestErr
}

type AuthenticateRequest struct {
	UserID string `json:"user_id"`
}

type IntrospectRequest struct {
	AccessToken string `json:"access_token"`
}

type Authenticate struct {
	AccessToken string `json:"access_token"`
}

type AuthenticateConcurrent struct {
	Result *Authenticate
	Error  *resterrors.RestErr
}
