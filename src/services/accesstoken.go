package services

import (
	"sync"

	"github.com/sshindanai/bookstore-utils-go/resterrors"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/models"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/repository/db"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/repository/rest"
)

type Service interface {
	GetAccessTokenByUserID(*models.AuthenticateRequest, chan *models.AuthenticateConcurrent)
	Introspection(*models.IntrospectRequest, chan *models.AccessTokenConcurrent)
	Create(*models.AccessTokenRequest, chan *models.AccessTokenConcurrent)
	Refresh(*models.IntrospectRequest, chan *models.AccessTokenConcurrent)
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

func (s *service) GetAccessTokenByUserID(req *models.AuthenticateRequest, output chan *models.AuthenticateConcurrent) {
	if err := req.Validate(); err != nil {
		res := &models.AuthenticateConcurrent{
			Error: err,
		}
		output <- res
		return
	}

	go s.repository.GetAccessTokenByUserID(req, output)
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

func (s *service) Introspection(at *models.IntrospectRequest, output chan *models.AccessTokenConcurrent) {
	if at.AccessToken == "" {
		res := &models.AccessTokenConcurrent{
			Error: resterrors.NewUnauthorizedError("invalid access token"),
		}
		output <- res
		return
	}

	go s.repository.Introspection(at, output)
}

func (s *service) Refresh(req *models.IntrospectRequest, output chan *models.AccessTokenConcurrent) {
	userCh := make(chan *models.AccessTokenConcurrent)
	go s.repository.Introspection(req, userCh)
	user := <-userCh

	if user.Error != nil {
		res := &models.AccessTokenConcurrent{
			Error: user.Error,
		}
		output <- res
		return
	}
	var wg sync.WaitGroup
	var err *resterrors.RestErr
	wg.Add(1)
	go func() {
		err = s.repository.DeleteKey(req.AccessToken)
		wg.Done()
	}()
	wg.Wait()
	if err != nil {
		res := &models.AccessTokenConcurrent{
			Error: err,
		}
		output <- res
		return
	}

	// generate new access token
	at := models.NewAccessToken(user.Result.UserID)
	at.Generate()
	go s.repository.Create(&at, output)
}
