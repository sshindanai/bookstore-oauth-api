package http

import (
	"net/http"

	"github.com/sshindanai/repo/bookstore-oauth-api/src/models"

	"github.com/labstack/echo"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/services"
)

type AccessTokenHandler interface {
	GetAccessTokenByUserID(echo.Context) error
	Introspection(echo.Context) error
	Create(echo.Context) error
	Refresh(echo.Context) error
}

type accessTokenHandler struct {
	service services.Service
}

func NewHandler(service services.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (handler *accessTokenHandler) GetAccessTokenByUserID(c echo.Context) error {
	id := c.Param("userid")
	req := models.AuthenticateRequest{
		UserID: id,
	}
	output := make(chan *models.AuthenticateConcurrent)

	go handler.service.GetAccessTokenByUserID(&req, output)
	res := <-output

	if res.Error != nil {
		return c.JSON(res.Error.Code, res.Error)
	}

	return c.JSON(http.StatusOK, res.Result)
}

// Get User info by access token
func (handler *accessTokenHandler) Introspection(c echo.Context) error {
	var request models.IntrospectRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid json body")
	}
	output := make(chan *models.AccessTokenConcurrent)

	go handler.service.Introspection(&request, output)
	res := <-output

	if res.Error != nil {
		return c.JSON(res.Error.Code, res.Error)
	}

	return c.JSON(http.StatusOK, res.Result)
}

// Create access token after authenticate
func (handler *accessTokenHandler) Create(c echo.Context) error {
	var request models.AccessTokenRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid json body")
	}
	output := make(chan *models.AccessTokenConcurrent)

	go handler.service.Create(&request, output)
	result := <-output

	if result.Error != nil {
		return c.JSON(result.Error.Code, result.Error)
	}

	return c.JSON(http.StatusOK, result.Result)
}

// Refresh access token while token is not expired
func (handler *accessTokenHandler) Refresh(c echo.Context) error {
	var request models.IntrospectRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid json body")
	}
	output := make(chan *models.AccessTokenConcurrent)

	go handler.service.Refresh(&request, output)
	res := <-output
	if res.Error != nil {
		return c.JSON(res.Error.Code, res.Error)
	}
	return c.JSON(http.StatusOK, res.Result)
}
