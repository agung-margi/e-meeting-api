package handler

import (
	"e-meeting-api/internal/domain/repository"
	"e-meeting-api/internal/usecase"
	"e-meeting-api/pkg/database"

	"github.com/labstack/echo/v4"
)

func RoutingRestAPI(e *echo.Echo) error {
	db, err := database.NewPostgresConnection("localhost", "postgres", "Bekasi2021$", "dbmeeting", "5432")
	if err != nil {
		e.Logger.Error(err)
		return err
	}

	// defer db.Close()
	// Inisialisasi repositori
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewStudentUseCase(userRepo)
	userHandler := NewUserHandler(userUsecase)

	//Routing
	e.GET("/users/:id", userHandler.GetUser)
	e.POST("/users", userHandler.SaveUser)
	e.PUT("/users/:id", userHandler.UpdateUser)
	return nil

}
