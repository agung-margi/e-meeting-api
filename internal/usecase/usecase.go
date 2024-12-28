package usecase

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/presenter/model"
	"time"
)

type UserUseCase interface {
	GetUser(ctx context.Context, id int) (*entity.User, error)
	SaveUser(ctx context.Context, user *model.UserRequest) error
	UpdateUser(ctx context.Context, id int, user *model.UpdateUserRequest) error
	Login(ctx context.Context, username string, password string) (string, *entity.User, error)
}

type ReservationUseCase interface {
	GetReservation(ctx context.Context, id int) ([]model.ReservationDetailsResponse, error)
	SaveReservation(ctx context.Context, reservationRequest *model.ReservationRequest, userId int) (*model.ReservationResponse, error)
	CheckAvailability(ctx context.Context, roomId int, startTime string, endTime string) (bool, error)
	GetReservationsByRoomAndDate(ctx context.Context, roomId int, date string) ([]entity.RoomSchedule, error)
	GetAll(ctx context.Context, startDate *time.Time, endDate *time.Time, roomType int, status string, userID *int) ([]*entity.Reservation, error)
	PayReservation(ctx context.Context, id int, userId int, isAdmin bool) error
	CancelReservation(ctx context.Context, id int, userId int, isAdmin bool) error
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
	GetByID(ctx context.Context, id int) (*entity.RoomWithType, error)
	SaveRoom(ctx context.Context, room *model.RoomRequest) error
	UpdateRoom(ctx context.Context, id int, room *model.RoomRequest) error
	DeleteRoom(ctx context.Context, id int) error
}

type InquiryUsecase interface {
	Save(ctx context.Context, reservation *entity.Reservation) (*entity.Inquiry, error)
}
