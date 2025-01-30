package repository

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/presenter/model"
	"time"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	SaveUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, id int, user *entity.User) error
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	SaveResetToken(ctx context.Context, email string, token string, expiredAt time.Time) error
	DeleteResetToken(ctx context.Context, email string) error
	UpdatePassword(ctx context.Context, id int, password string) error
	ValidateResetToken(ctx context.Context, email string, token string) (bool, error)
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
	GetAll(ctx context.Context, name string, roomType int, capacity int) ([]entity.Room, error)
	GetByID(ctx context.Context, id int) (*entity.Room, error)
	SaveRoom(ctx context.Context, room *entity.Room) error
	UpdateRoom(ctx context.Context, id int, room *entity.Room) error
	DeleteRoom(ctx context.Context, id int) error
	RoomExists(ctx context.Context, name string) (bool, error)
}

type ReservationRepository interface {
	GetByID(ctx context.Context, id int) (*entity.Reservation, error)
	Save(ctx context.Context, reservation *entity.Reservation, details *entity.ReservationDetails) (*entity.Reservation, error)
	CheckAvailability(ctx context.Context, roomId int, startTime string, endTime string) (bool, error)
	GetReservationsByRoomAndDate(ctx context.Context, roomId int, date string) ([]entity.RoomSchedule, error)
	GetRoomPriceByID(ctx context.Context, roomID int) (int, error)
	GetSnackPriceByID(ctx context.Context, snackID int) (int, error)
	GetReservationDetails(ctx context.Context, reservationID int) (*[]model.ReservationDetailsResponse, error)
	GetAll(ctx context.Context, startDate, endDate string, roomType int, status string, userID *int) ([]*entity.ReservationHistory, error)
	UpdateStatus(ctx context.Context, id int, status string) error
	GetReservationsCountByRoomAndDate(ctx context.Context, roomID int, date string) (int, error)
	GetSchedulesByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entity.RoomSchedule, error)
	GetDashboardDataByDateRange(ctx context.Context, startDate, endDate time.Time) (entity.Dashboard, error)
	GetRoomOmzetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entity.RoomOmzetDetails, error)
	GetExpiredReservations(ctx context.Context, now time.Time) ([]entity.Reservation, error)
}

type InquiryRepository interface {
	Save(context.Context, *entity.Inquiry) (*entity.Inquiry, error)
	GetByID(ctx context.Context, id int, userId int) (*entity.Inquiry, error)
	DeleteInquiry(ctx context.Context, id int, userId int) error
}
