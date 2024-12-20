package usecase

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/internal/domain/repository"
	"e-meeting-api/presenter/model"
	"fmt"
	"log"
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

func (u *reservationUseCase) SaveReservation(ctx context.Context, reservationRequest *model.ReservationRequest) (*model.ReservationResponse, error) {

	startTime, err := time.Parse(time.RFC3339, reservationRequest.StartTime)
	if err != nil {
		log.Printf("Error parsing start time: %v", err)
		return nil, err
	}
	endTime, err := time.Parse(time.RFC3339, reservationRequest.EndTime)
	if err != nil {
		log.Printf("Error parsing end time: %v", err)
		return nil, err
	}

	roomPricePerHour, err := u.repo.GetRoomPriceByID(ctx, reservationRequest.RoomID)
	if err != nil {
		log.Printf("Error retrieving room price: %v", err)
		return nil, err
	}

	duration := endTime.Sub(startTime).Hours()
	if duration <= 0 {
		return nil, fmt.Errorf("invalid reservation duration: start time should be before end time")
	}
	totalRoomPrice := int(duration) * roomPricePerHour

	snackPrice, err := u.repo.GetSnackPriceByID(ctx, reservationRequest.SnackID)
	if err != nil {
		log.Printf("Error retrieving snack price: %v", err)
		return nil, err
	}

	totalSnackPrice := snackPrice * reservationRequest.Participants
	totalPrice := totalRoomPrice + totalSnackPrice

	reservation := &entity.Reservation{
		UserID:      reservationRequest.UserID,
		RoomID:      reservationRequest.RoomID,
		StartTime:   startTime,
		EndTime:     endTime,
		BookingDate: reservationRequest.BookingDate,
		RoomPrice:   totalRoomPrice,
		SnackPrice:  totalSnackPrice,
		TotalPrice:  totalPrice,
		Status:      "booked",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = u.repo.SaveReservation(ctx, reservation, &entity.ReservationDetails{
		Name:         reservationRequest.Name,
		Phone:        reservationRequest.Phone,
		Company:      reservationRequest.Company,
		SnackID:      reservationRequest.SnackID,
		Participants: reservationRequest.Participants,
		Notes:        reservationRequest.Notes,
	})
	if err != nil {
		log.Printf("Error saving reservation: %v", err)
		return nil, err
	}

	details, err := u.repo.GetReservationDetails(ctx, reservation.ID)
	if err != nil {
		log.Printf("Error retrieving reservation details: %v", err)
		return nil, err
	}

	response := &model.ReservationResponse{
		ReservationID: reservation.ID,
		UserID:        reservation.UserID,
		RoomID:        reservation.RoomID,
		StartTime:     reservation.StartTime,
		EndTime:       reservation.EndTime,
		BookingDate:   reservation.BookingDate,
		RoomPrice:     reservation.RoomPrice,
		SnackPrice:    reservation.SnackPrice,
		TotalPrice:    reservation.TotalPrice,
		Status:        reservation.Status,
		Details:       details,
	}

	return response, nil
}
