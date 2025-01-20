package handler

import (
	"context"
	"e-meeting-api/configs"
	"e-meeting-api/internal/domain/repository"
	"e-meeting-api/internal/usecase"
	"e-meeting-api/pkg/database"
	"e-meeting-api/pkg/middleware"

	_ "e-meeting-api/presenter/handler/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func RoutingRestAPI(g *echo.Group, logger echo.Logger, ctx context.Context) error {

	configs.LoadConfig()

	DBHOST := configs.AppConfig.Database.Host
	DBUSERNAME := configs.AppConfig.Database.User
	DBPASSWORD := configs.AppConfig.Database.Password
	DBNAME := configs.AppConfig.Database.DBName
	DBPORT := configs.AppConfig.Database.Port

	db, err := database.NewPostgresConnection(DBHOST, DBUSERNAME, DBPASSWORD, DBNAME, DBPORT)
	if err != nil {
		logger.Error(err)
		return err
	}

	// Inisialisasi repositori dan usecase
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewStudentUseCase(userRepo)
	userHandler := NewUserHandler(userUsecase)

	g.GET("/swagger/*", echoSwagger.WrapHandler)
	g.POST("/login", userHandler.Login)

	// Routing API User
	g.POST("/register", userHandler.SaveUser)
	g.GET("/users/:id", userHandler.GetUser, middleware.AuthMiddleware)
	g.PUT("/users/:id", userHandler.UpdateUser, middleware.AuthMiddleware)
	g.POST("/forget-password", userHandler.ForgetPassword)
	g.POST("/reset-password/:token", userHandler.ResetPassword)
	// Routing API Snack
	snackRepo := repository.NewSnackRepository(db)
	snackUsecase := usecase.NewSnackUseCase(snackRepo)
	snackHandler := NewSnackHandler(snackUsecase)

	g.GET("/snacks", snackHandler.GetSnacks, middleware.AuthMiddleware)

	// Routing API RoomType
	roomtypeRepo := repository.NewRoomtypeRepository(db)
	roomtypeUsecase := usecase.NewRoomtypeUseCase(roomtypeRepo)
	roomtypeHandler := NewRoomTypeHandler(roomtypeUsecase)

	g.GET("/roomtypes", roomtypeHandler.GetRoomTypes, middleware.AuthMiddleware)

	// Routing API Room
	roomRepo := repository.NewRoomRepository(db)
	roomUsecase := usecase.NewRoomUseCase(roomRepo)
	roomHandler := NewRoomHandler(roomUsecase)

	g.GET("/rooms", roomHandler.GetRooms, middleware.AuthMiddleware)
	g.GET("/rooms/:id", roomHandler.GetRoom, middleware.AuthMiddleware)
	g.POST("/rooms", roomHandler.SaveRoom, middleware.AuthMiddleware, middleware.IsAdminMiddleware)
	g.PUT("/rooms/:id", roomHandler.UpdateRoom, middleware.AuthMiddleware, middleware.IsAdminMiddleware)
	g.DELETE("/rooms/:id", roomHandler.DeleteRoom, middleware.AuthMiddleware, middleware.IsAdminMiddleware)

	// Routing API Inquiry & Reservation
	inquiryRepo := repository.NewInquiryRepository(db)
	reservationRepo := repository.NewReservationRepository(db)
	inquiryUsecase := usecase.NewInquiryUseCase(inquiryRepo, roomRepo, snackRepo, reservationRepo)
	inquiryHandler := NewInquiryHandler(inquiryUsecase)
	reservationUsecase := usecase.NewReservationUseCase(reservationRepo, inquiryRepo)
	reservationHandler := NewReservationHandler(reservationUsecase)

	// Start expired reservation worker
	go func() {
		reservationUsecase.StartExpiredReservationWorker(ctx)
	}()

	g.POST("/reservations/inquiry", inquiryHandler.SaveInquiry, middleware.AuthMiddleware)
	g.POST("/reservations/:inquiry_id", reservationHandler.SaveReservation, middleware.AuthMiddleware)
	g.GET("/reservations/:id", reservationHandler.GetReservation, middleware.AuthMiddleware)
	g.GET("/reservations", reservationHandler.GetReservations, middleware.AuthMiddleware)
	g.PUT("/reservations/:id/pay", reservationHandler.PayReservation, middleware.AuthMiddleware)
	g.PUT("/reservations/:id/cancel", reservationHandler.CancelReservation, middleware.AuthMiddleware)
	g.GET("/room-schedule", reservationHandler.GetRoomSchedule, middleware.AuthMiddleware)
	g.GET("/reservations/schedules", reservationHandler.GetSchedules, middleware.AuthMiddleware)
	g.GET("/dashboard", reservationHandler.GetDashboardData, middleware.AuthMiddleware, middleware.IsAdminMiddleware)
	//

	return nil
}
