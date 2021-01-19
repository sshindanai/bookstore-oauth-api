package rest

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mercadolibre/golang-restclient/rest"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer() /*I commented flagParse() in func init() of the rest package to make it works*/
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	// Arrange
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8082/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: -1,
		RespBody:     `{}`,
	})

	repository := userRepository{}

	// Act
	user, err := repository.LoginUser("email@gmail.com", "password")

	// Assert
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Code)
	assert.EqualValues(t, "invalid restclient response when trying to login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	// Arrange
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8082/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: http.StatusUnauthorized,
		RespBody:     `{"message": "invalid login credentials", "status_code": "401", "error": "unauthorized"}`,
	})

	repository := userRepository{}

	// Act
	user, err := repository.LoginUser("email@gmail.com", "password")

	// Assert
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Code)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)
}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	// Arrange
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8082/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: http.StatusUnauthorized,
		RespBody:     `{"message": "invalid login credentials", "status_code": 401, "error": "unauthorized"}`,
	})

	repository := userRepository{}

	// Act
	user, err := repository.LoginUser("email@gmail.com", "password")

	// Assert
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Code)
	assert.EqualValues(t, "invalid login credentials", err.Message)
}

func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	// Arrange
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8082/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": "123", "first_name": "test123", "last_name": "test456", "email": "email@gmail.com"}`,
	})

	repository := userRepository{}

	// Act
	user, err := repository.LoginUser("email@gmail.com", "password")

	// Assert
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Code)
	assert.EqualValues(t, "error when trying to unmarshal users login response", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	// Arrange
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8082/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": 123, "first_name": "test123", "last_name": "test456", "email": "email@gmail.com"}`,
	})

	repository := userRepository{}

	// Act
	user, err := repository.LoginUser("email@gmail.com", "password")

	// Assert
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.EqualValues(t, 123, user.Id)
	assert.EqualValues(t, "test123", user.FirstName)
	assert.EqualValues(t, "test456", user.LastName)
	assert.EqualValues(t, "email@gmail.com", user.Email)
}
