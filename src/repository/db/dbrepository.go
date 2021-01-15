package db

import (
	"context"
	"strconv"
	"time"

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
	GetById(string, chan *models.AccessTokenConcurrent)
	Create(*models.AccessToken, chan *models.AccessTokenConcurrent)
}
type dbRepository struct{}

func (db *dbRepository) GetById(id string, output chan *models.AccessTokenConcurrent) {
	go func() {
		rdb := redis.NewRedis()
		token, err := rdb.Get(ctx, id).Result()
		if err != nil || token == "" {
			output <- &models.AccessTokenConcurrent{
				Error: errors.ParseError(err),
			}
			return
		}

		userIDStr := rdb.HGet(ctx, token, "user_id").Val()
		userID, _ := strconv.ParseInt(userIDStr, 10, 64)

		clientIDStr := rdb.HGet(ctx, token, "client_id").Val()
		clientID, _ := strconv.ParseInt(clientIDStr, 10, 64)

		expStr := rdb.HGet(ctx, token, "expires").Val()
		exp, _ := strconv.ParseInt(expStr, 10, 64)

		result := &models.AccessToken{
			AccessToken: token,
			UserID:      userID,
			ClientID:    clientID,
			Expires:     exp,
		}
		output <- &models.AccessTokenConcurrent{
			Result: result,
		}
		return
	}()
}

func (db *dbRepository) Create(at *models.AccessToken, output chan *models.AccessTokenConcurrent) {
	go func() {
		rdb := redis.NewRedis()

		_, err := rdb.SetEX(ctx, strconv.Itoa(int(at.UserID)), at.AccessToken, time.Hour*expirationTime).Result()
		if err != nil {
			output <- &models.AccessTokenConcurrent{
				Error: errors.NewUnauthorizedError(err.Error()),
			}
			return
		}

		_, err = rdb.HSet(ctx, at.AccessToken, "user_id", at.UserID, "client_id", at.ClientID, "expires", at.Expires).Result()
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
