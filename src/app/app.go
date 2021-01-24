package app

import (
	"github.com/labstack/echo"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/http"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/repository/db"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/repository/rest"
	"github.com/sshindanai/repo/bookstore-oauth-api/src/services"
)

var router *echo.Echo

func init() {
	router = echo.New()
}

func StartApp() {
	atHandler := http.NewHandler(services.NewService(rest.NewRestUsersRepository(), db.NewRepository()))
	endpointsRegister(atHandler)
	router.Logger.Fatal(router.Start(":1323"))
}

func endpointsRegister(handler http.AccessTokenHandler) {
	router.GET("/oauth/accesstoken/:userid", handler.GetAccessTokenByUserID)
	router.GET("/oauth/accesstoken/introspect", handler.Introspection)
	router.POST("/oauth/access_token", handler.Create)
	router.PUT("/oauth/accesstoken/refresh", handler.Refresh)
}
