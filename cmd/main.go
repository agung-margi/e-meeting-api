package main

import (
	"context"
	"e-meeting-api/configs"
	"e-meeting-api/presenter/handler"
	_ "e-meeting-api/presenter/handler/docs"

	"github.com/labstack/echo/v4"
	m "github.com/labstack/echo/v4/middleware"
)

// @title E-Meeting API
// @version 1.0
// @description Ini adalah API untuk E-Meeting
// @termsOfService http://swagger.io/terms/
// @host localhost:8080
// @BasePath /api/v1

// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name Authorization

func main() {
	configs.LoadConfig()
	e := echo.New()
	e.Use(m.CORSWithConfig(m.CORSConfig{
		AllowOrigins: []string{"*"}, // Mengizinkan semua asal (origin)
		AllowMethods: []string{"*"}, // Mengizinkan semua metode HTTP
		AllowHeaders: []string{"*"}, // Mengizinkan semua header
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e.Static("/photos", configs.AppConfig.PhotosPath)
	v1 := e.Group("/api/v1")
	if err := handler.RoutingRestAPI(v1, e.Logger, ctx); err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
