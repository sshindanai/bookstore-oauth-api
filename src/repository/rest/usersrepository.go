package rest

import (
	"encoding/json"
	"time"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/sshindanai/bookstore-utils-go/resterrors"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/models"
)

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8080",
		Timeout: 100 * time.Millisecond,
	}
)

type RestUsersRepository interface {
	LoginUser(string, string) (*models.User, *resterrors.RestErr)
}

type userRepository struct{}

func NewRestUsersRepository() RestUsersRepository {
	return &userRepository{}
}

func (r *userRepository) LoginUser(email string, password string) (*models.User, *resterrors.RestErr) {
	request := models.UserLoginRequest{
		Email:    email,
		Password: password,
	}

	response := usersRestClient.Post("/users/login", request)
	if response == nil || response.Response == nil {
		return nil, resterrors.NewUnauthorizedError("invalid restclient response when trying to login user")
	}

	if response.StatusCode > 299 {
		var restErr resterrors.RestErr
		if err := json.Unmarshal(response.Bytes(), &restErr); err != nil {
			// Happen when tag "status_code" is string
			return nil, resterrors.NewInternalServerError("invalid error interface when trying to login user", err)
		}
		return nil, &restErr
	}

	var user models.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, resterrors.NewUnauthorizedError("error when trying to unmarshal users login response")
	}
	return &user, nil
}
