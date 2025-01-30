package usecase

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/presenter/model"
	"time"

	"github.com/xuri/excelize/v2"
)

type UserUseCase interface {
	GetUser(ctx context.Context, id int) (*entity.User, error)
	SaveUser(ctx context.Context, user *model.UserRequest) error
	UpdateUser(ctx context.Context, id int, user *model.UpdateUserRequest) error
	Login(ctx context.Context, username string, password string) (string, *entity.User, error)
	GetResetPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, email string, token string, newPassword string) error
}

type ReservationUseCase interface {
	GetReservation(ctx context.Context, reservationID int) (*[]model.ReservationDetailsResponse, error)
	Save(ctx context.Context, inquiryId int, userId int) (*[]model.ReservationDetailsResponse, error)
	GetReservationsByRoomAndDate(ctx context.Context, roomId int, date string) ([]entity.RoomSchedule, error)
	GetAll(ctx context.Context, startDate string, endDate string, roomType int, status string, userID *int) ([]*entity.ReservationHistory, error)
	PayReservation(ctx context.Context, id int, userId int, isAdmin bool) error
	CancelReservation(ctx context.Context, id int, userId int, isAdmin bool) error
	GetSchedulesByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entity.RoomSchedule, error)
	GetDashboardData(ctx context.Context, startDate, endDate time.Time) (entity.Dashboard, error)
	ExpireReservations(ctx context.Context) error
	StartExpiredReservationWorker(ctx context.Context)
	ExportReservations(ctx context.Context, startDate string, endDate string, roomType int, status string, userID *int) (*excelize.File, error)
}

type SnackUseCase interface {
	GetAll(ctx context.Context) ([]entity.Snack, error)
	GetByID(ctx context.Context, id int) (*entity.Snack, error)
}

type RoomtypeUseCase interface {
	GetAll(ctx context.Context) ([]entity.RoomtypeRepository, error)
	GetByID(ctx context.Context, id int) (*entity.RoomtypeRepository, error)
	GetByName(ctx context.Context, name string) (*entity.RoomtypeRepository, error)
}

type RoomUseCase interface {
	GetAll(ctx context.Context, name string, roomType int, capacity int) ([]entity.Room, error)
	GetByID(ctx context.Context, id int) (*entity.Room, error)
	SaveRoom(ctx context.Context, room *model.RoomRequest) error
	UpdateRoom(ctx context.Context, id int, room *model.RoomRequest) error
	DeleteRoom(ctx context.Context, id int) error
	GetByName(ctx context.Context, name string) (bool, error)
}

type InquiryUsecase interface {
	Save(ctx context.Context, reservation *entity.Reservation) (*entity.Inquiry, error)
	GetByID(ctx context.Context, id int, userId int) (*entity.Inquiry, error)
}
