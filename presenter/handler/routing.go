package handler

import (
	"e-meeting-api/internal/domain/repository"
	"e-meeting-api/internal/usecase"
	"e-meeting-api/pkg/database"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func RoutingRestAPI(e *echo.Echo) error {
	err := godotenv.Load(".env")
	if err != nil {
		e.Logger.Error(err)
	}

	DBHOST := os.Getenv("DB_HOST")
	DBUSERNAME := os.Getenv("DB_USERNAME")
	DBPASSWORD := os.Getenv("DB_PASSWORD")
	DBNAME := os.Getenv("DB_NAME")
	DBPORT := os.Getenv("DB_PORT")

	db, err := database.NewPostgresConnection(DBHOST, DBUSERNAME, DBPASSWORD, DBNAME, DBPORT)
	if err != nil {
		e.Logger.Error(err)
		return err
	}

	// Inisialisasi repositori dan usecase
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewStudentUseCase(userRepo)
	userHandler := NewUserHandler(userUsecase)

	// Routing API
	e.GET("/users/:id", userHandler.GetUser)
	e.POST("/users", userHandler.SaveUser)
	e.PUT("/users/:id", userHandler.UpdateUser)
	e.POST("/login", userHandler.Login)
	e.POST("/logout", userHandler.Logout)

	return nil
}
