package handler

import (
	"e-meeting-api/internal/domain/repository"
	"e-meeting-api/internal/usecase"
	"e-meeting-api/pkg/database"
	"e-meeting-api/pkg/middleware"
	"os"

	_ "e-meeting-api/presenter/handler/docs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
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

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/login", userHandler.Login)
	e.POST("/logout", userHandler.Logout)

	// Routing API User
	e.POST("/users", userHandler.SaveUser)
	e.GET("/users/:id", userHandler.GetUser, middleware.AuthMiddleware)
	e.PUT("/users/:id", userHandler.UpdateUser, middleware.AuthMiddleware)

	// Routing API Snack
	snackRepo := repository.NewSnackRepository(db)
	snackUsecase := usecase.NewSnackUseCase(snackRepo)
	snackHandler := NewSnackHandler(snackUsecase)

	e.GET("/snacks", snackHandler.GetSnacks, middleware.AuthMiddleware)

	// Routing API RoomType
	roomtypeRepo := repository.NewRoomtypeRepository(db)
	roomtypeUsecase := usecase.NewRoomtypeUseCase(roomtypeRepo)
	roomtypeHandler := NewRoomTypeHandler(roomtypeUsecase)

	e.GET("/roomtypes", roomtypeHandler.GetRoomTypes, middleware.AuthMiddleware)

	// Routing API Room
	roomRepo := repository.NewRoomRepository(db)
	roomUsecase := usecase.NewRoomUseCase(roomRepo)
	roomHandler := NewRoomHandler(roomUsecase)

	e.GET("/rooms", roomHandler.GetRooms, middleware.AuthMiddleware)
	e.GET("/rooms/:id", roomHandler.GetRoom, middleware.AuthMiddleware)
	e.POST("/rooms", roomHandler.SaveRoom, middleware.AuthMiddleware, middleware.IsAdminMiddleware)
	e.PUT("/rooms/:id", roomHandler.UpdateRoom, middleware.AuthMiddleware, middleware.IsAdminMiddleware)

	// Routing API Reservation
	reservationRepo := repository.NewReservationRepository(db)
	reservationUsecase := usecase.NewReservationUseCase(reservationRepo)
	reservationHandler := NewReservationHandler(reservationUsecase)

	e.POST("/reservations", reservationHandler.SaveReservation, middleware.AuthMiddleware)
	e.GET("/reservations/:id", reservationHandler.GetReservation, middleware.AuthMiddleware)
	e.GET("/reservations", reservationHandler.GetReservations, middleware.AuthMiddleware)
	e.PUT("/reservations/:id/pay", reservationHandler.PayReservation, middleware.AuthMiddleware)
	e.PUT("/reservations/:id/cancel", reservationHandler.CancelReservation, middleware.AuthMiddleware)
	// e.GET("/get-availability", reservationHandler.CheckAvailability)

	return nil
}

// 	//Routing API Reservation

// 	e.POST("/reservations", reservationHandler.SaveReservation)
// 	e.GET("/get-availability", reservationHandler.CheckAvailability)
// 	e.GET("/get-room-schedule", reservationHandler.CheckRoomSchedule)
