package repository

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/presenter/model"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	SaveUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, id int, user *entity.User) error
	DeleteUser(ctx context.Context, id string) error
	CheckEmailExists(ctx context.Context, email string) (bool, error)
}

type SnackRepository interface {
	GetAll(ctx context.Context) ([]entity.Snack, error)
	GetByID(ctx context.Context, id int) (*entity.Snack, error)
}

type RoomtypeRepository interface {
	GetAll(ctx context.Context) ([]entity.RoomtypeRepository, error)
	GetByID(ctx context.Context, id int) (*entity.RoomtypeRepository, error)
	GetByName(ctx context.Context, name string) (*entity.RoomtypeRepository, error)
	SaveRoomType(ctx context.Context, roomType *entity.RoomtypeRepository) error
}

type RoomRepository interface {
	GetAll(ctx context.Context) ([]entity.Room, error)
	GetByID(ctx context.Context, id int) (*entity.RoomWithType, error)
	SaveRoom(ctx context.Context, room *entity.Room) error
	UpdateRoom(ctx context.Context, id int, room *entity.Room) error
}

type ReservationRepository interface {
	GetByID(ctx context.Context, id int) (*entity.Reservation, error)
	SaveReservation(ctx context.Context, reservation *entity.Reservation, details *entity.ReservationDetails) error
	CheckAvailability(ctx context.Context, roomId int, startTime string, endTime string) (bool, error)
	GetReservationsByRoomAndDate(ctx context.Context, roomId int, date string) ([]*entity.Reservation, error)
	GetRoomPriceByID(ctx context.Context, roomID int) (int, error)
	GetSnackPriceByID(ctx context.Context, snackID int) (int, error)
	GetReservationDetails(ctx context.Context, reservationID int) ([]model.ReservationDetailsResponse, error)
	GetAll(ctx context.Context) ([]*entity.Reservation, error)
	GetByUserID(ctx context.Context, userID int) ([]*entity.Reservation, error)
	UpdateStatus(ctx context.Context, id int, status string) error
}
