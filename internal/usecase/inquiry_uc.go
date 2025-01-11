package usecase

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/internal/domain/repository"
	"errors"
	"fmt"
	"math"
	"time"
)

type inquiryUsecase struct {
	inquiryRepo     repository.InquiryRepository
	roomRepo        repository.RoomRepository
	snackRepo       repository.SnackRepository
	reservationRepo repository.ReservationRepository
}

func NewInquiryUseCase(inquiryRepo repository.InquiryRepository, roomRepo repository.RoomRepository, snackRepo repository.SnackRepository, reservationRepo repository.ReservationRepository) InquiryUsecase {
	return &inquiryUsecase{inquiryRepo: inquiryRepo, roomRepo: roomRepo, snackRepo: snackRepo, reservationRepo: reservationRepo}
}

func (u *inquiryUsecase) Save(ctx context.Context, reservation *entity.Reservation) (*entity.Inquiry, error) {

	// Get room details
	room, err := u.roomRepo.GetByID(ctx, reservation.RoomID)
	if err != nil {
		return nil, err
	}

	// Get snack details (if any)
	var snack *entity.Snack
	if reservation.SnackID != nil {

		snack, err = u.snackRepo.GetByID(ctx, *reservation.SnackID)
		if err != nil {
			return nil, err
		}
	}

	// Parse start and end time
	date := reservation.BookingDate
	start := reservation.StartTime
	end := reservation.EndTime
	format := "2006-01-02 15:04"

	startStr := fmt.Sprintf("%s %s", date, start)
	endStr := fmt.Sprintf("%s %s", date, end)

	startTime, err := time.Parse(format, startStr)
	if err != nil {
		return nil, errors.New("invalid start time format")
	}

	endTime, err := time.Parse(format, endStr)
	if err != nil {
		return nil, errors.New("invalid end time format")
	}

	// Validation checks
	if startTime.After(endTime) {
		return nil, errors.New("start time should be before end time")
	}

	if startTime.Before(time.Now()) {
		return nil, errors.New("start time should be in the future")
	}

	if reservation.Participants < 1 {
		return nil, errors.New("participants should be at least 1")
	}

	if reservation.Participants > room.Capacity {
		return nil, errors.New("participants should not exceed room capacity")
	}

	// Check room availability
	available, err := u.reservationRepo.CheckAvailability(ctx, reservation.RoomID, startStr, endStr)
	if err != nil {
		return nil, errors.New("failed to check room availability")
	}

	if !available {
		return nil, errors.New("room is not available for the selected time")
	}

	// Calculate room price
	duration := int(math.Ceil(endTime.Sub(startTime).Hours()))
	if duration < 1 {
		duration = 1 // Minimum duration of 1 hour
	}
	roomTotalPrice := room.Price * duration

	// Calculate snack price
	snackTotalPrice := 0
	if snack != nil {
		snackTotalPrice = snack.Price * reservation.Participants
	}

	// Calculate grand total
	grandTotal := roomTotalPrice + snackTotalPrice

	// Create inquiry
	inquiry := &entity.Inquiry{
		UserID:          reservation.UserID,
		RoomID:          reservation.RoomID,
		RoomName:        room.Name,
		RoomType:        room.RoomType,
		RoomCapacity:    room.Capacity,
		RoomPrice:       room.Price,
		BookingDate:     reservation.BookingDate,
		StartTime:       startStr,
		EndTime:         endStr,
		Duration:        duration,
		Name:            reservation.Name,
		Phone:           reservation.Phone,
		Company:         reservation.Company,
		Participants:    reservation.Participants,
		SnackID:         nil,
		SnackName:       "",
		SnackCategory:   "",
		SnackPrice:      0,
		TotalRoomPrice:  roomTotalPrice,
		TotalSnackPrice: snackTotalPrice,
		TotalPrice:      grandTotal,
		Notes:           reservation.Notes,
		CreatedAt:       time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:       time.Now().Format("2006-01-02 15:04:05"),
	}
	// Update inquiry with snack details if available
	if snack != nil {
		inquiry.SnackPrice = snack.Price
		inquiry.SnackID = reservation.SnackID
		inquiry.SnackName = snack.Name
		inquiry.SnackCategory = snack.Category
		inquiry.TotalSnackPrice = snackTotalPrice
		inquiry.TotalPrice = grandTotal
	}

	// Save inquiry
	inquiry, err = u.inquiryRepo.Save(ctx, inquiry)
	if err != nil {
		return nil, err
	}
	return inquiry, nil
}

func (u *inquiryUsecase) GetByID(ctx context.Context, id int, userId int) (*entity.Inquiry, error) {
	return u.inquiryRepo.GetByID(ctx, id, userId)
}
