package usecase

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/internal/domain/repository"
	"e-meeting-api/presenter/model"
	"errors"
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

func (u *reservationUseCase) GetReservationsByRoomAndDate(ctx context.Context, roomId int, date string) ([]entity.RoomSchedule, error) {
	return u.repo.GetReservationsByRoomAndDate(ctx, roomId, date)
}
func (u *reservationUseCase) PayReservation(ctx context.Context, id int, userId int, isAdmin bool) error {

	reservation, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if reservation.UserID != userId && !isAdmin {
		return errors.New("unauthorized")
	}

	// validasi status pembayaran
	if reservation.Status == "paid" {
		return errors.New("reservation that has already been paid")
	}

	if time.Now().After(reservation.StartTime) {
		return errors.New("reservation that has already started")
	}

	if time.Now().After(reservation.EndTime) {
		return errors.New("reservation that has already ended")
	}

	if reservation.Status == "cancelled" {
		return errors.New("reservation that has already been cancelled")
	}

	return u.repo.UpdateStatus(ctx, id, "paid")
}

func (u *reservationUseCase) CancelReservation(ctx context.Context, id int, userId int, isAdmin bool) error {
	reservation, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if reservation.Status == "paid" {
		return errors.New("reservation that has already been paid")
	}

	if time.Now().After(reservation.StartTime) {
		return errors.New("reservation that has already started")
	}

	if time.Now().After(reservation.EndTime) {
		return errors.New("reservation that has already ended")
	}

	if reservation.Status == "cancelled" {
		return errors.New("reservation that has already been cancelled")
	}

	return u.repo.UpdateStatus(ctx, id, "cancelled")
}

func (u *reservationUseCase) SaveReservation(ctx context.Context, reservationRequest *model.ReservationRequest, userId int) (*model.ReservationResponse, error) {

	startTime, err := time.Parse(time.RFC3339, reservationRequest.StartTime)
	if err != nil {
		return nil, err
	}
	endTime, err := time.Parse(time.RFC3339, reservationRequest.EndTime)
	if err != nil {
		return nil, err
	}

	if startTime.After(endTime) {
		return nil, errors.New("invalid reservation duration: start time should be before end time")
	}

	if startTime.Before(time.Now()) {
		return nil, errors.New("invalid reservation start time: start time should be in the future")
	}

	avalilable, err := u.CheckAvailability(ctx, reservationRequest.RoomID, reservationRequest.StartTime, reservationRequest.EndTime)
	if err != nil {
		return nil, err
	}

	if !avalilable {
		return nil, errors.New("room is not available for the selected date and time")
	}

	roomPricePerHour, err := u.repo.GetRoomPriceByID(ctx, reservationRequest.RoomID)
	if err != nil {
		return nil, err
	}

	duration := endTime.Sub(startTime).Hours()
	if duration <= 0 {
		return nil, errors.New("invalid reservation duration: start time should be before end time")
	}
	totalRoomPrice := int(duration) * roomPricePerHour

	snackPrice := 0

	if reservationRequest.SnackID != nil {
		snackPrice, err = u.repo.GetSnackPriceByID(ctx, *reservationRequest.SnackID)
		if err != nil {
			return nil, err
		}
	}

	snackID := reservationRequest.SnackID
	totalSnackPrice := snackPrice * reservationRequest.Participants
	totalPrice := totalRoomPrice + totalSnackPrice

	reservation := &entity.Reservation{
		UserID:      userId,
		RoomID:      reservationRequest.RoomID,
		StartTime:   startTime,
		EndTime:     endTime,
		BookingDate: time.Now(),
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
		SnackID:      snackID,
		Participants: reservationRequest.Participants,
		Notes:        reservationRequest.Notes,
	})
	if err != nil {
		return nil, err
	}

	details, err := u.repo.GetReservationDetails(ctx, reservation.ID)
	if err != nil {
		return nil, err
	}

	response := &model.ReservationResponse{
		ReservationID: reservation.ID,
		UserID:        reservation.UserID,
		RoomID:        reservation.RoomID,
		StartTime:     reservation.StartTime,
		EndTime:       reservation.EndTime,
		BookingDate:   time.Now(),
		RoomPrice:     reservation.RoomPrice,
		SnackPrice:    reservation.SnackPrice,
		TotalPrice:    reservation.TotalPrice,
		Status:        reservation.Status,
		Details:       details,
	}

	return response, nil
}

func (u *reservationUseCase) GetAll(ctx context.Context, startDate *time.Time, endDate *time.Time, roomType int, status string, userID *int) ([]*entity.Reservation, error) {
	reservations, err := u.repo.GetAll(ctx, startDate, endDate, roomType, status, userID)
	if err != nil {
		return nil, err
	}

	return reservations, nil
}
