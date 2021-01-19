package services

import (
	"github.com/sshindanai/repo/bookstore-oauth-api/src/models"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/repository/db"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/repository/rest"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/utils/errors"
)

type Service interface {
	GetByID(string, chan *models.AuthenticateConcurrent)
	Introspection(string, chan *models.AccessTokenConcurrent)
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

func (s *service) GetByID(id string, output chan *models.AuthenticateConcurrent) {
	go s.repository.GetByID(id, output)
}

func (s *service) Create(request *models.AccessTokenRequest, output chan *models.AccessTokenConcurrent) {
	if err := request.Validate(); err != nil {
		res := &models.AccessTokenConcurrent{
			Error: err,
		}
		output <- res
		return
	}

	// TODO: Support both grant types: client_credentials and password
	// Authenticate the user against the Users API:

	user, err := s.restUserRepo.LoginUser(request.Email, request.Password)
	if err != nil {
		res := &models.AccessTokenConcurrent{
			Error: err,
		}
		output <- res
		return
	}

	// generate a new access token
	at := models.NewAccessToken(user.Id)
	at.Generate()

	// Save the new access token to redis
	go s.repository.Create(&at, output)
}

func (s *service) Introspection(at string, output chan *models.AccessTokenConcurrent) {
	if at == "" {
		res := &models.AccessTokenConcurrent{
			Error: errors.NewUnauthorizedError("invalid access token"),
		}
		output <- res
		return
	}

	go s.repository.Introspection(at, output)
}

func (s *service) Refresh(token string, output chan *models.AccessTokenConcurrent) {
	userCh := make(chan *models.AccessTokenConcurrent)
	go s.repository.Introspection(token, userCh)
	user := <-userCh

	if user.Error != nil {
		res := &models.AccessTokenConcurrent{
			Error: errors.NewUnauthorizedError("access token has expired already"),
		}
		output <- res
		return
	}

	// generate new access token
	at := models.NewAccessToken(user.Result.UserID)
	at.Generate()
	go s.repository.Create(&at, output)
}
