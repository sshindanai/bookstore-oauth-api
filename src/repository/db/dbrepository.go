package db

import (
	"context"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/sshindanai/bookstore-utils-go/resterrors"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/clients/redis"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/models"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/utils/errors"
)

var ctx = context.Background()

const (
	expirationTime = 24
)

func NewRepository() Repository {
	return &dbRepository{}
}

type Repository interface {
	GetAccessTokenByUserID(*models.AuthenticateRequest, chan *models.AuthenticateConcurrent)
	Introspection(*models.IntrospectRequest, chan *models.AccessTokenConcurrent)
	Create(*models.AccessToken, chan *models.AccessTokenConcurrent)
	DeleteKey(string) *resterrors.RestErr
}

type dbRepository struct{}

func (db *dbRepository) GetAccessTokenByUserID(req *models.AuthenticateRequest, output chan *models.AuthenticateConcurrent) {
	go func() {
		redisClient := redis.NewRedis()
		token, err := redisClient.Get(ctx, req.UserID).Result()
		if err != nil || token == "" {
			output <- &models.AuthenticateConcurrent{
				Error: errors.ParseError(err, req.UserID),
			}
			return
		}

		res := &models.Authenticate{
			AccessToken: token,
		}

		output <- &models.AuthenticateConcurrent{
			Result: res,
		}
		return
	}()
}

func (db *dbRepository) Introspection(at *models.IntrospectRequest, output chan *models.AccessTokenConcurrent) {
	go func() {
		redisClient := redis.NewRedis()
		data, err := redisClient.HGetAll(ctx, at.AccessToken).Result()
		if err != nil || data == nil {
			output <- &models.AccessTokenConcurrent{
				Error: resterrors.NewUnauthorizedError(err.Error()),
			}
			return
		}

		var userInfo models.AccessToken
		err = mapstructure.WeakDecode(data, &userInfo)
		if err != nil {
			output <- &models.AccessTokenConcurrent{
				Error: resterrors.NewUnauthorizedError(err.Error()),
			}
			return
		}

		if userInfo.UserID < 1 {
			output <- &models.AccessTokenConcurrent{
				Error: resterrors.NewUnauthorizedError("access token doesn't exist"),
			}
			return
		}
		result := &models.AccessTokenConcurrent{
			Result: &userInfo,
		}
		output <- result
	}()
}

func (db *dbRepository) Create(at *models.AccessToken, output chan *models.AccessTokenConcurrent) {
	go func() {
		redisClient := redis.NewRedis()
		_, err := redisClient.SetEX(ctx, strconv.Itoa(int(at.UserID)), at.AccessToken, time.Hour*expirationTime).Result()
		if err != nil {
			output <- &models.AccessTokenConcurrent{
				Error: resterrors.NewUnauthorizedError(err.Error()),
			}
			return
		}

		_, err = redisClient.HSet(ctx, at.AccessToken, "UserID", at.UserID, "ClientId", at.ClientID, "Expires", at.Expires).Result()
		if err != nil {
			output <- &models.AccessTokenConcurrent{
				Error: resterrors.NewUnauthorizedError(err.Error()),
			}
			return
		}
		res := &models.AccessTokenConcurrent{
			Result: at,
		}
		output <- res
		return
	}()
}

func (db *dbRepository) DeleteKey(key string) *resterrors.RestErr {
	redisClient := redis.NewRedis()

	_, err := redisClient.Del(ctx, key).Result()
	if err != nil {
		return resterrors.NewBadRequestError(err.Error())
	}
	return nil
}
