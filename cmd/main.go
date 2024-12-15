package main

import (
	"e-meeting-api/presenter/handler"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	if err := handler.RoutingRestAPI(e); err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Fatal(e.Start(":8080"))
}
