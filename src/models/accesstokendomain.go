package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/sshindanai/bookstore-utils-go/resterrors"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/utils/cryptoutils"
)

const (
	expirationTime            = 24
	grantTypePassword         = "password"
	grantTypeClientCredential = "client_credentials"
)

func NewAccessToken(userID int64) AccessToken {
	return AccessToken{
		UserID:   userID,
		Expires:  time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
		ClientID: 456,
	}
}

func (a *AccessToken) IsExpired() *resterrors.RestErr {
	now := time.Now().UTC()
	expirationTime := time.Unix(a.Expires, 0)

	if now.After(expirationTime) {
		return resterrors.NewUnauthorizedError("access token is expired")
	}

	return nil
}

func (at *AccessToken) Generate() {
	at.AccessToken = cryptoutils.GetSHA256(fmt.Sprintf("at-%d-%d-ran", at.UserID, at.Expires))
}

func (at *AccessToken) Validate() *resterrors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return resterrors.NewBadRequestError("invalid access token id")
	}
	if at.UserID <= 0 {
		return resterrors.NewBadRequestError("invalid user id")
	}
	if at.ClientID <= 0 {
		return resterrors.NewBadRequestError("invalid client id")
	}
	if at.Expires <= 0 {
		return resterrors.NewBadRequestError("invalid expiration time")
	}
	return nil
}

func (at *AccessTokenRequest) Validate() *resterrors.RestErr {
	switch at.GrantType {
	case grantTypePassword,
		grantTypeClientCredential:
		break
	default:
		return resterrors.NewUnauthorizedError("invalid grant_type parameter")
	}

	return nil
}

func (t *AuthenticateRequest) Validate() *resterrors.RestErr {
	t.UserID = strings.TrimSpace(t.UserID)
	if t.UserID <= "" {
		return resterrors.NewBadRequestError("invalid user id")
	}
	return nil
}
