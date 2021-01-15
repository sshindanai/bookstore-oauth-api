package rest

import (
	"encoding/json"
	"time"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/models"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/utils/errors"
)

var (
	userRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8082",
		Timeout: 100 * time.Millisecond,
	}
)

type RestUsersRepository interface {
	LoginUser(string, string) (*models.User, *errors.RestErr)
}

type userRepository struct{}

func NewRestUsersRepository() RestUsersRepository {
	return &userRepository{}
}

func (r *userRepository) LoginUser(email string, password string) (*models.User, *errors.RestErr) {
	request := models.UserLoginRequest{
		Email:    email,
		Password: password,
	}

	response := userRestClient.Post("/users/login", request)
	if response == nil || response.Response == nil {
		return nil, errors.NewUnauthorizedError("invalid restclient response when trying to login user")
	}

	// if response.StatusCode > 299 {
	// 	apiErr, err := rest_errors.NewRestErrorFromBytes(response.Bytes())
	// 	if err != nil {
	// 		return nil, rest_errors.NewInternalServerError("invalid error interface when trying to login user", err)
	// 	}
	// 	return nil, apiErr
	// }

	var user models.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, errors.NewUnauthorizedError("error when trying to unmarshal users login response")
	}
	return &user, nil
}
