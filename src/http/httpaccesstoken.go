package http

import (
	"net/http"

	"github.com/sshindanai/repo/bookstore-oauth-api/src/models"

	"github.com/labstack/echo"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/services"
)

type AccessTokenHandler interface {
	GetById(echo.Context) error
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

func (handler *accessTokenHandler) GetById(c echo.Context) error {
	accessTokenID := c.Param("accesstokenid")

	output := make(chan *models.AccessTokenConcurrent)

	go handler.service.GetById(accessTokenID, output)
	res := <-output

	if res.Error != nil {
		return c.JSON(res.Error.Code, res.Error)
	}

	return c.JSON(http.StatusOK, res.Result)
}

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

func (handler *accessTokenHandler) Refresh(c echo.Context) error {
	userID := c.Param("accesstokenid")
	output := make(chan *models.AccessTokenConcurrent)

	go handler.service.Refresh(userID, output)
	res := <-output
	if res.Error != nil {
		return c.JSON(res.Error.Code, res.Error)
	}

	return c.JSON(http.StatusOK, res.Result)
}
