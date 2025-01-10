package main

import (
	"context"
	"e-meeting-api/presenter/handler"
	_ "e-meeting-api/presenter/handler/docs"

	"github.com/labstack/echo/v4"
)

// @title E-Meeting API
// @version 1.0
// @description Ini adalah API untuk E-Meeting
// @termsOfService http://swagger.io/terms/
// @host localhost:8080
// @BasePath /

// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name Authorization

func main() {
	e := echo.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := handler.RoutingRestAPI(e, ctx); err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
