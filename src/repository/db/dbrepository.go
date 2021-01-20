package db

import (
	"context"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
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
	GetByID(string, chan *models.AuthenticateConcurrent)
	Introspection(string, chan *models.AccessTokenConcurrent)
	Create(*models.AccessToken, chan *models.AccessTokenConcurrent)
}
type dbRepository struct{}

func (db *dbRepository) GetByID(id string, output chan *models.AuthenticateConcurrent) {
	go func() {
		redisClient := redis.NewRedis()
		token, err := redisClient.Get(ctx, id).Result()
		if err != nil || token == "" {
			output <- &models.AuthenticateConcurrent{
				Error: errors.ParseError(err, id),
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

func (db *dbRepository) Introspection(at string, output chan *models.AccessTokenConcurrent) {
	go func() {
		redisClient := redis.NewRedis()
		data, err := redisClient.HGetAll(ctx, at).Result()
		if err != nil || data == nil {
			output <- &models.AccessTokenConcurrent{
				Error: errors.NewUnauthorizedError(err.Error()),
			}
			return
		}

		var userInfo models.AccessToken
		err = mapstructure.WeakDecode(data, &userInfo)
		if err != nil {
			output <- &models.AccessTokenConcurrent{
				Error: errors.NewUnauthorizedError(err.Error()),
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
				Error: errors.NewUnauthorizedError(err.Error()),
			}
			return
		}

		_, err = redisClient.HSet(ctx, at.AccessToken, "UserId", at.UserID, "ClientId", at.ClientID, "Expires", at.Expires).Result()
		if err != nil {
			output <- &models.AccessTokenConcurrent{
				Error: errors.NewUnauthorizedError(err.Error()),
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
