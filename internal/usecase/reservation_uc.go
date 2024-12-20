package usecase

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/internal/domain/repository"
	"e-meeting-api/presenter/model"
	"errors"
	"fmt"
	"time"
)

type reservationUseCase struct {
	repo repository.ReservationRepository
}

func NewReservationUseCase(repo repository.ReservationRepository) ReservationUseCase {
	return &reservationUseCase{repo: repo}
}

func (u *reservationUseCase) GetReservation(ctx context.Context, id int) (*entity.Reservation, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *reservationUseCase) SaveReservation(ctx context.Context, reservation *model.ReservationRequest) (*entity.Reservation, error) {

	startTime, err := time.Parse(time.RFC3339, reservation.StartTime)
	if err != nil {
		return nil, errors.New("Invalid start time format")
	}

	endTime, err := time.Parse(time.RFC3339, reservation.EndTime)
	if err != nil {
		return nil, errors.New("Invalid start time format")
	}

	newReservation := entity.Reservation{
		UserID:      reservation.UserID,
		RoomID:      reservation.RoomID,
		StartTime:   startTime,
		EndTime:     endTime,
		BookingDate: reservation.BookingDate,
		RoomPrice:   reservation.RoomPrice,
		SnackPrice:  reservation.SnackPrice,
		TotalPrice:  reservation.TotalPrice,
		Status:      reservation.Status,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
	fmt.Println(newReservation)

	err = u.repo.SaveReservation(ctx, &newReservation)
	if err != nil {
		return &newReservation, err
	}
	return &newReservation, err
}

func (u *reservationUseCase) CheckAvailability(ctx context.Context, roomId int, startTime string, endTime string) (bool, error) {
	avalilable, err := u.repo.CheckAvailability(ctx, roomId, startTime, endTime)
	if err != nil {
		return false, err
	}

	return !avalilable, nil
}

func (u *reservationUseCase) GetReservationsByRoomAndDate(ctx context.Context, roomId int, date string) ([]*entity.Reservation, error) {
	return u.repo.GetReservationsByRoomAndDate(ctx, roomId, date)
}
