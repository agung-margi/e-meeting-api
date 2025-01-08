package main

import (
	"e-meeting-api/presenter/handler"
	_ "e-meeting-api/presenter/handler/docs"

	"github.com/labstack/echo/v4"
)

// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name Authorization

func main() {

	e := echo.New()
	if err := handler.RoutingRestAPI(e); err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Fatal(e.Start(":8080"))
}