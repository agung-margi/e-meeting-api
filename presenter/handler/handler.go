package handler

import "github.com/labstack/echo/v4"

type UserHandler interface {
	GetUser(c echo.Context) error
	SaveUser(c echo.Context) error
	UpdateUser(c echo.Context) error
}
