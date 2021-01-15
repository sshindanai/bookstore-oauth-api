package services

import (
	"sync"

	"github.com/sshindanai/repo/bookstore-oauth-api/src/models"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/repository/db"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/repository/rest"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/utils/errors"
)

type Service interface {
	GetById(string, chan *models.AccessTokenConcurrent)
	Create(*models.AccessTokenRequest, chan *models.AccessTokenConcurrent)
	Refresh(string, chan *models.AccessTokenConcurrent)
}

type service struct {
	restUserRepo rest.RestUsersRepository
	repository   db.Repository
}

func NewService(restUserRepo rest.RestUsersRepository, repo db.Repository) Service {
	return &service{
		restUserRepo: restUserRepo,
		repository:   repo,
	}
}

func (s *service) GetById(id string, output chan *models.AccessTokenConcurrent) {
	go s.repository.GetById(id, output)
}

func (s *service) Create(request *models.AccessTokenRequest, output chan *models.AccessTokenConcurrent) {
	if err := request.Validate(); err != nil {
		res := &models.AccessTokenConcurrent{
			Error: err,
		}
		output <- res
		return
	}

	// TODO: Support both grant types
	// user, err := s.restUserRepo.LoginUser(request.Username, request.Password)
	// if err != nil {
	// 	return nil, err
	// }

	// Mock login
	var wg sync.WaitGroup
	wg.Add(1)
	user, err := func(email string, password string) (*models.User, *errors.RestErr) {
		wg.Done()
		return &models.User{
			Id:        123,
			FirstName: "Shindanai",
			LastName:  "Mongkolsin",
			Email:     "sshindanai.m@gmail.com",
		}, nil
	}(request.Username, request.Password)
	if err != nil {
		res := &models.AccessTokenConcurrent{
			Error: err,
		}
		output <- res
		return
	}
	wg.Wait()

	// generate new access token
	at := models.NewAccessToken(user.Id)
	at.Generate()

	// Save the new access token to redis
	go s.repository.Create(&at, output)
}

func (s *service) Refresh(userID string, output chan *models.AccessTokenConcurrent) {
	userCh := make(chan *models.AccessTokenConcurrent)
	go s.repository.GetById(userID, userCh)
	user := <-userCh

	if user.Error != nil {
		res := &models.AccessTokenConcurrent{
			Error: user.Error,
		}
		output <- res
		return
	}

	if user.Result.IsExpired() {
		res := &models.AccessTokenConcurrent{
			Error: errors.NewUnauthorizedError("token is already expired"),
		}
		output <- res
		return
	}

	// generate new access token
	at := models.NewAccessToken(user.Result.UserID)
	at.Generate()
	go s.repository.Create(&at, output)
}
