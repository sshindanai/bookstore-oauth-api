package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/sshindanai/repo/bookstore-oauth-api/src/utils/cryptoutils"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/utils/errors"
)

func NewAccessToken(userID int64) AccessToken {
	return AccessToken{
		UserID:   userID,
		Expires:  time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
		ClientID: 456,
	}
}

func (a *AccessToken) IsExpired() bool {
	now := time.Now().UTC()
	expirationTime := time.Unix(a.Expires, 0)

	return now.After(expirationTime)
}

func (at *AccessToken) Generate() {
	at.AccessToken = cryptoutils.GetSHA256(fmt.Sprintf("at-%d-%d-ran", at.UserID, at.Expires))
}

func (at *AccessToken) Validate() *errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return errors.NewBadRequestError("invalid access token id")
	}
	if at.UserID <= 0 {
		return errors.NewBadRequestError("invalid user id")
	}
	if at.ClientID <= 0 {
		return errors.NewBadRequestError("invalid client id")
	}
	if at.Expires <= 0 {
		return errors.NewBadRequestError("invalid expiration time")
	}
	return nil
}

func (at *AccessTokenRequest) Validate() *errors.RestErr {
	switch at.GrantType {
	case grantTypePassword:
		break
	case grantTypeClientCredential:
		break
	default:
		return errors.NewUnauthorizedError("invalid grant_type parameter")
	}

	return nil
}
